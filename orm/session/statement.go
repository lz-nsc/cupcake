package session

import "strings"

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
)

// Statement is used to construct a complete SQL statement with
// given value and action
type statement struct {
	sql     map[Type]string
	sqlVars map[Type][]interface{}
}

func (s *statement) Set(action Type, args ...interface{}) {
	if s.sql == nil {
		s.sql = make(map[Type]string)
		s.sqlVars = make(map[Type][]interface{})
	}
	clause, vars := generators[action](args...)
	s.sql[action] = clause
	s.sqlVars[action] = vars
}

func (s *statement) Build(actions ...Type) (string, []interface{}) {
	var sqls []string
	var vars []interface{}
	for _, action := range actions {
		if sql, ok := s.sql[action]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, s.sqlVars[action]...)
		}
	}

	return strings.Join(sqls, " "), vars
}
