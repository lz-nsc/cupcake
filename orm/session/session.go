package session

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"

	"github.com/lz-nsc/cupcake/orm/log"
	"github.com/lz-nsc/cupcake/orm/schema"
	"github.com/lz-nsc/cupcake/orm/translator"
)

type Session struct {
	db        *sql.DB
	sql       strings.Builder
	sqlVars   []interface{}
	trans     translator.Translator
	schema    *schema.Schema
	statement *statement
}

func New(db *sql.DB, trans translator.Translator) *Session {
	return &Session{
		db:        db,
		trans:     trans,
		statement: &statement{},
	}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
}

func (s Session) DB() *sql.DB {
	return s.db
}

func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

func (s *Session) Exec() (result sql.Result, err error) {
	// Reset session after successfully execute previous query
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) QueryRow() *sql.Row {
	// Reset session after successfully execute previous query
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	// Reset session after successfully execute previous query
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

// Insert generate proper insert statement with given structs
// and execute to database:
//	INSERT INTO table_name (column_list)
//	VALUES
//		(value_list_1),
//		(value_list_2),
//		...
//		(value_list_n);
func (s *Session) Insert(instances ...interface{}) (int64, error) {
	var modelName string
	values := make([]interface{}, 0)
	for _, instance := range instances {
		err := s.Model(instance)
		if err != nil {
			return 0, err
		}
		schema := s.Schema()
		if modelName == "" {
			modelName = schema.Name
		}
		if schema.Name != modelName {
			// values should be mapped to same table
			return 0, errors.New("different model instances passed to insert")
		}

		values = append(values, schema.GetValues(instance))

	}
	if schema := s.Schema(); schema != nil {
		s.statement.Set(INSERT, schema.Name, schema.FieldNames)
	}
	s.statement.Set(VALUES, values...)

	sql, vars := s.statement.Build(INSERT, VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// List accept array of struct as parameter and return all the record of related table
func (s *Session) List(arr interface{}) error {
	list := reflect.Indirect(reflect.ValueOf(arr))
	elemType := list.Type().Elem()

	// Generate a new instance with given list's element type, then pass to orm session
	err := s.Model(reflect.New(elemType).Elem().Interface())
	if err != nil {
		return err
	}
	schema := s.Schema()
	s.statement.Set(SELECT, schema.Name, schema.FieldNames)
	sql, vars := s.statement.Build(SELECT)

	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	for rows.Next() {
		elem := reflect.New(elemType).Elem()

		var values []interface{}
		for _, name := range schema.FieldNames {
			values = append(values, elem.FieldByName(name).Addr().Interface())
		}

		if err := rows.Scan(values...); err != nil {
			return err
		}

		list.Set(reflect.Append(list, elem))
	}
	return nil
}
