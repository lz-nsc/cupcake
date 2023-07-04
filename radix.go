package cupcake

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type nodeType uint8
type methodType uint8

const (
	StaticNode nodeType = iota
	RegrexNode
	ParamNode
	WildNode
)

const (
	GET methodType = iota
	POST
	PUT
	DELETE
)

var methodMapping = map[string]methodType{
	http.MethodGet:    GET,
	http.MethodPost:   POST,
	http.MethodPut:    PUT,
	http.MethodDelete: DELETE,
}

var (
	ErrNotAllow = errors.New("Method Not Allowed")
	ErrNotFound = errors.New("NOT FOUND")
)

type radixNode struct {
	head      byte
	tail      byte
	nodeType  nodeType
	prefix    string
	rex       *regexp.Regexp
	children  []radixNodes
	endpoints map[methodType]*endpoint
}
type radixNodes []*radixNode

type endpoint struct {
	paramKeys []string
	handler   HandlerFunc
}

func NewNode(prefix string) *radixNode {
	node := &radixNode{
		prefix: prefix,

		children: make([]radixNodes, 4),
	}
	if len(prefix) > 0 {
		node.head = prefix[0]
	}
	return node
}

func (node *radixNode) InsertNode(path string, method methodType, handler HandlerFunc) {
	curNode := node
	search := path
	paramKeys := []string{}
	for {
		if len(search) == 0 {
			curNode.setEndpoint(paramKeys, method, handler)
			return
		}
		pType, pattern, _, _, _, _, _ := parsePath(search)
		next := curNode.findNext(pType, search)

		if (pType == ParamNode || pType == RegrexNode) && pattern != "" {
			paramKeys = append(paramKeys, pattern)
		}

		// next node not found
		if next == nil {
			child, nextStart := curNode.addNode(search)
			search = search[nextStart:]
			curNode = child
			continue
		}

		commonPrefix := longestCommonPrefix(search, next.prefix)
		search = search[commonPrefix:]

		// Search share prefix with current node, the search the rest part amoung children
		if commonPrefix == len(next.prefix) {
			curNode = next
			continue
		}

		// Slip node
		child := &radixNode{
			head:     next.prefix[commonPrefix],
			nodeType: next.nodeType,
			prefix:   next.prefix[commonPrefix:],
		}

		next.prefix = next.prefix[:commonPrefix]

		child.nodeType = next.nodeType
		child.children = next.children
		child.endpoints = next.endpoints

		// add new childNode to parent node
		next.children = make([]radixNodes, 4)
		next.children[child.nodeType] = []*radixNode{child}

		curNode = next
	}
}

func (node *radixNode) Route(path string, method methodType) (handler HandlerFunc, params map[string]string, wild string, err error) {
	params = map[string]string{}
	child, paramVals, wildStr, routeErr := node.route(path, method)
	if routeErr != nil {
		err = routeErr
		return
	}

	endpoint := child.endpoints[method]
	if len(endpoint.paramKeys) != len(paramVals) {
		panic("ParamKeys and ParamVals do not match")
	}
	for idx := 0; idx < len(endpoint.paramKeys); idx++ {
		params[endpoint.paramKeys[idx]] = paramVals[idx]
	}
	handler = endpoint.handler
	wild = wildStr

	return
}

func (node *radixNode) route(path string, method methodType) (child *radixNode, paramVals []string, wild string, err error) {
	curNode := node
	search := path
	for t, nodeGroup := range curNode.children {
		var tempCurrent *radixNode
		tempSearch := search
		nType := nodeType(t)
		switch nType {
		case StaticNode:
			next := curNode.findNext(StaticNode, path)
			if next == nil || !strings.HasPrefix(path, next.prefix) {
				continue
			}
			tempCurrent = next
			tempSearch = tempSearch[len(next.prefix):]
		case ParamNode, RegrexNode:
			for _, next := range nodeGroup {
				tempCurrent = next
				tailIdx := strings.IndexByte(tempSearch, next.tail)
				if tailIdx < 0 {
					if next.tail == '/' {
						tailIdx = len(tempSearch)
					} else {
						continue
					}
				}
				if (nType == RegrexNode && tailIdx == 0) || strings.IndexByte(tempSearch[:tailIdx], '/') != -1 {
					continue
				}

				if nType == RegrexNode {
					if tempCurrent.rex == nil {
						panic("Regex node with empty rex")
					}

					if !tempCurrent.rex.MatchString(tempSearch[:tailIdx]) {
						continue
					}
				}

				valsSize := len(paramVals)
				if tempCurrent.prefix != "" {
					paramVals = append(paramVals, tempSearch[:tailIdx])
				}
				tempSearch = tempSearch[tailIdx:]
				if len(tempSearch) == 0 {
					// Find endpoints
					if tempCurrent.endpoints != nil && tempCurrent.endpoints[method] != nil {
						child = tempCurrent
						err = nil
						return
					}
					err = ErrNotAllow
					continue
				}

				// Keep searching
				res, vals, wildStr, newErr := tempCurrent.route(tempSearch, method)
				if res != nil {
					child = res
					paramVals = append(paramVals, vals...)
					wild = wildStr
					err = nil
					return
				}
				if err != ErrNotAllow {
					err = newErr
				}
				// Fail to find leaf, recover paramVals
				paramVals = paramVals[:valsSize]
				tempSearch = search
			}
		default:
			if len(nodeGroup) > 0 {
				tempCurrent = nodeGroup[0]
				wild = tempSearch
				tempSearch = ""
			}
		}

		if tempCurrent == nil {
			continue
		}
		//Found node
		if len(tempSearch) == 0 {
			if tempCurrent.endpoints != nil && tempCurrent.endpoints[method] != nil {
				child = tempCurrent
				err = nil
				return
			}
			err = ErrNotAllow
			continue
		}

		res, vals, wildStr, newErr := tempCurrent.route(tempSearch, method)
		// Found node
		if res != nil {
			child = res
			paramVals = append(paramVals, vals...)
			wild = wildStr
			err = nil
			return
		}
		// Fail to find node in this group
		if err != ErrNotFound {
			err = newErr
		}
	}

	if err == nil {
		err = ErrNotFound
	}
	return
}

func (node *radixNode) addNode(path string) (*radixNode, int) {
	pType, pattern, regex, tail, _, _, nextStart := parsePath(path)
	child := NewNode(pattern)
	child.tail = tail
	child.nodeType = pType
	if node.children[pType] == nil {
		node.children[pType] = []*radixNode{}
	}

	if regex != "" {
		rex, err := regexp.Compile(regex)
		if err != nil {
			panic(fmt.Sprintf("Invalid regexp pattern '%s' in path", regex))
		}
		child.rex = rex
	}

	node.children[pType] = append(node.children[pType], child)
	return child, nextStart
}

// Parse Path recognize note type from path, if the path is special type,
// then return the sepecial part and it's start and end index
func parsePath(path string) (
	nType nodeType, pattern string, regexPattern string, tail byte,
	start int, end int, nextStart int,
) {
	// Default pattern is the path itself
	nType = StaticNode
	pattern = path
	start = 0
	end = len(path)
	nextStart = len(path)
	tail = '/'

	paramIdx := strings.Index(path, "{")
	wildIdx := strings.Index(path, "*")
	if paramIdx < 0 && wildIdx < 0 {
		return
	}

	if paramIdx < 0 && wildIdx >= 0 && wildIdx < len(pattern)-1 {
		panic("Wildcard '*' must be at the end of the path")
	}

	// Contains param or regex
	if paramIdx >= 0 {
		if paramIdx > 0 {
			// Segment the static part
			pattern = pattern[start:paramIdx]
			end = paramIdx
			nextStart = paramIdx
			return
		}
		// Remove '{'
		start = 1
		// Search for '}'
		endParamIdx := strings.Index(path, "}")
		if endParamIdx < 0 || endParamIdx < paramIdx {
			panic("Invalid path")
		}

		end = endParamIdx
		pattern = pattern[start:end]
		if next := strings.Index(pattern, "{"); next >= 0 {
			panic("Nested parantheses is not allowed")
		}

		// Remove '}'
		nextStart = endParamIdx + 1

		// Start with param or regex
		if paramIdx == 0 {
			regexIdx := strings.Index(pattern, ":")
			if regexIdx >= 0 {
				if regexIdx == end-1 {
					panic("Invalid regex path")
				}
				// Regex
				nType = RegrexNode
				regexPattern = pattern[regexIdx+1:]
				pattern = pattern[:regexIdx]
				if nextStart < len(path) {
					tail = path[nextStart]
				}
				return
			}

			// Param
			if nextStart < len(path) && path[nextStart] != '/' {
				panic("Params in path must followed with '/'")
			}
			nType = ParamNode
			return
		}
	}

	// Wildcast
	if wildIdx > 0 {
		// Segment the static part
		end = wildIdx
		nextStart = wildIdx
		pattern = pattern[start:end]
		return
	}

	nType = WildNode
	pattern = "*"
	return
}
func (node *radixNode) findNext(pType nodeType, pattern string) *radixNode {
	if pType != StaticNode {
		return nil
	}

	if pattern == "" || node.children[pType] == nil {
		return nil
	}
	// Find next node
	for _, next := range node.children[pType] {
		if next.head != pattern[0] {
			continue
		}
		return next
	}
	return nil
}

func (node *radixNode) setEndpoint(paramKeys []string, method methodType, handler HandlerFunc) {
	if node.endpoints == nil {
		node.endpoints = make(map[methodType]*endpoint)
	}

	node.endpoints[method] = &endpoint{
		paramKeys: paramKeys,
		handler:   handler,
	}
}

func longestCommonPrefix(origin string, target string) int {
	length := len(origin)
	if len(target) < length {
		length = len(target)
	}
	for idx := 0; idx < length; idx++ {
		if origin[idx] != target[idx] {
			return idx
		}
	}
	return length
}

func parseMethod(method string) methodType {
	if m, ok := methodMapping[method]; ok {
		return m
	}
	panic(fmt.Sprintf("Invalid mathod: %s", method))
}
