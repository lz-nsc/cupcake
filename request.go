package cupcake

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/lz-nsc/cupcake/log"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
)

type Request struct {
	req    *http.Request
	path   string
	method string
	params map[string]string
	data   []byte
	wild   string
}

const (
	ApplicationJSON  = "application/json"
	ApplicationXML   = "application/xml"
	ApplicationForm  = "application/x-www-form-urlencoded"
	MultipartForm    = "multipart/form-data"
	ApplicationProto = "application/x-protobuf"
	ApplicationYAML  = "application/x-yaml"
	TextXML          = "text/xml"
)

func NewRequest(r *http.Request) *Request {
	req := &Request{
		req:    r,
		path:   r.URL.Path,
		method: r.Method,
		params: make(map[string]string),
		data:   []byte{},
	}

	return req
}
func (r *Request) PostForm(key string) string {
	return r.req.FormValue(key)
}

func (r *Request) Query(key string) string {
	return r.req.URL.Query().Get(key)
}

func (r Request) String() string {
	return fmt.Sprintf("Request: %s %s", r.method, r.path)
}

func (r Request) Path() string {
	return r.path
}

func (r Request) Method() string {
	return r.method
}
func (r *Request) SetParams(params map[string]string) {
	r.params = params
}
func (r *Request) SetWild(wild string) {
	r.wild = wild
}

func (r Request) Params() map[string]string {
	return r.params
}
func (r Request) Param(key string) string {
	return r.params[key]
}

func (r Request) Wild() string {
	return r.wild
}

func (r Request) Body() io.ReadCloser {
	return r.req.Body
}

func (r Request) Parse(obj interface{}) error {
	ct := r.req.Header.Get("Content-Type")
	if len(ct) == 0 {
		return r.parseJson(obj)
	}
	vals := strings.Split(ct, ";")
	switch vals[0] {
	case ApplicationForm:
		return r.parseForm(obj)
	case ApplicationProto:
		return r.parseProto(obj)
	case ApplicationYAML:
		return r.parseYaml(obj)
	case ApplicationXML, TextXML:
		return r.parseXml(obj)
	case ApplicationJSON:
		return r.parseJson(obj)
	default:
		return errors.New("Unsupported Content-Type:" + ct)
	}
}

func (r *Request) readData() {
	if r.req.Body == nil {
		return
	}
	var data []byte
	// Decode gzip request
	if ce := r.req.Header.Values("Content-Encoding"); len(ce) != 0 {
		vals := strings.Split(ce[0], ",")
		for _, val := range vals {
			if val == "gzip" {
				gReader, err := gzip.NewReader(r.Body())
				if err != nil {
					return
				}
				data, _ = ioutil.ReadAll(gReader)
			}
		}
	}
	if data == nil {
		data, _ = ioutil.ReadAll(r.req.Body)
	}
	buffer := bytes.NewBuffer(data)
	r.req.Body = ioutil.NopCloser(buffer)
	r.data = data
}

func (r Request) parseJson(obj interface{}) error {
	return json.Unmarshal(r.data, obj)
}
func (r Request) parseForm(obj interface{}) error {
	err := r.req.ParseForm()
	if err != nil {
		return err
	}
	objType := reflect.TypeOf(obj)
	if objType.Kind() != reflect.Ptr && objType.Elem().Kind() != reflect.Struct {
		return errors.New("given obj is not pointer to a struct")
	}
	objValElm := reflect.ValueOf(obj).Elem()
	objTypeElm := objType.Elem()
	parseFormToStruct(r.req.Form, objTypeElm, objValElm)
	return nil
}

func parseFormToStruct(form url.Values, objType reflect.Type, objVal reflect.Value) {
	for idx := 0; idx < objType.NumField(); idx++ {
		field := objVal.Field(idx)
		if !field.CanSet() {
			continue
		}
		fieldType := objType.Field(idx)
		if fieldType.Anonymous && fieldType.Type.Kind() == reflect.Struct {
			parseFormToStruct(form, fieldType.Type, field)
			continue
		}
		fieldName, ok := getformName(fieldType)
		if !ok {
			continue
		}

		fieldVal := form[fieldName]
		if len(fieldVal) == 0 || fieldVal[0] == "" {
			continue
		}
		var err error
		// Convert string to the exact type of the field
		switch fieldType.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val, err := strconv.ParseInt(fieldVal[0], 10, 64)
			if err == nil {
				field.SetInt(val)
			}
		case reflect.Bool:
			val, err := strconv.ParseBool(fieldVal[0])
			if err == nil {
				field.SetBool(val)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val, err := strconv.ParseUint(fieldVal[0], 10, 64)
			if err == nil {
				field.SetUint(val)
			}
		case reflect.Float32, reflect.Float64:
			val, err := strconv.ParseFloat(fieldVal[0], 64)
			if err == nil {
				field.SetFloat(val)
			}
		case reflect.Interface:
			field.Set(reflect.ValueOf(fieldVal[0]))
		case reflect.String:
			field.SetString(fieldVal[0])
		case reflect.Struct:
			if fieldType.Type.String() == "time.Time" {
				t, err := parseFormTime(fieldVal[0])
				if err == nil {
					field.Set(reflect.ValueOf(t))
				}
			}
		case reflect.Slice:
			if fieldType.Type == reflect.TypeOf([]int(nil)) {
				formVals := form[fieldName]
				field.Set(reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(int(1))), len(formVals), len(formVals)))
				for i := 0; i < len(formVals); i++ {
					val, err := strconv.Atoi(formVals[i])
					if err == nil {
						field.Index(i).SetInt(int64(val))
					}
				}
			} else if fieldType.Type == reflect.TypeOf([]string(nil)) {
				formVals := form[fieldName]
				field.Set(reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf("")), len(formVals), len(formVals)))
				for i := 0; i < len(formVals); i++ {
					field.Index(i).SetString(formVals[i])
				}
			}
		}
		if err != nil {
			log.Debugf("parse bool error, err: %s", err.Error())
		}
	}
}

func (r Request) parseProto(obj interface{}) error {
	return proto.Unmarshal(r.data, obj.(proto.Message))
}
func (r Request) parseYaml(obj interface{}) error {
	return yaml.Unmarshal(r.data, obj)
}
func (r Request) parseXml(obj interface{}) error {
	return xml.Unmarshal(r.data, obj)
}
func getformName(field reflect.StructField) (string, bool) {
	tags := strings.Split(field.Tag.Get("form"), ",")
	var tag string
	if len(tags) == 0 || tags[0] == "" {
		// if no eplicit tag for this field, then use it name in the struct
		tag = field.Name
	} else if tags[0] == "-" {
		// skip this field
		return "", false
	} else {
		// user form tag as field name
		tag = tags[0]
	}

	return tag, true
}
func parseFormTime(timeStr string) (time.Time, error) {
	var pattern string
	if len(timeStr) >= 25 {
		timeStr = timeStr[:25]
		pattern = time.RFC3339
	} else if strings.HasSuffix(strings.ToUpper(timeStr), "Z") {
		pattern = time.RFC3339
	} else if len(timeStr) >= 19 {
		if strings.Contains(timeStr, "T") {
			pattern = "2006-01-02T15:04:05"
		} else {
			pattern = "2006-01-02 15:04:05"
		}
		timeStr = timeStr[:19]
	} else if len(timeStr) >= 10 {
		timeStr = timeStr[:10]
		pattern = "2006-01-02"
	} else if len(timeStr) >= 8 {
		timeStr = timeStr[:8]
		pattern = "15:04:05"
	}
	return time.ParseInLocation(pattern, timeStr, time.Local)
}
