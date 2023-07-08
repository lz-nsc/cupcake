package session

import (
	"fmt"
	"reflect"
	"strings"

	S "github.com/lz-nsc/cupcake/orm/schema"
)

// Model parse a struct in to schema and store it in the session
func (s *Session) Model(value interface{}) error {
	if s.schema == nil || reflect.TypeOf(value) != reflect.TypeOf(s.schema.Model) {
		schema, err := S.Parse(value, s.trans)
		if err != nil {
			return err
		}

		s.schema = schema
	}
	return nil
}

func (s Session) Schema() *S.Schema {
	return s.schema
}

// CreateTable use schema store in the session to generate and execute the CREATE TABLE SQL command
func (s Session) CreateTable() error {
	schema := s.Schema()
	if schema == nil {
		return fmt.Errorf("create table with empty schema")
	}
	var fields []string
	for _, field := range schema.FieldMap {
		fields = append(fields, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(fields, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", schema.Name, desc)).Exec()
	return err
}

// DropTable drop the table in database which is related to the schema stores in the session
func (s Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.Schema().Name)).Exec()
	return err
}

func (s Session) HasTable() bool {
	sql, values := s.trans.TableExistSQL(s.Schema().Name)
	row := s.Raw(sql, values...).QueryRow()
	var tmp interface{}

	err := row.Scan(&tmp)
	return err == nil
}
