package session

import (
	"errors"
	"reflect"
)

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

	// Build the complete statement
	sql, vars := s.statement.Build(SELECT, WHERE, ORDERBY, LIMIT)

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

// Update update records filter by WHERE clause with given map
// or key-value list("field1","value1","field2","value2"...)
func (s *Session) Update(data ...interface{}) (int64, error) {
	if s.Schema() == nil {
		return 0, errors.New("empty schema")
	}
	// If first parameter is map, then use this map
	values, ok := data[0].(map[string]interface{})
	if !ok {
		// If not, then construct a map with given parameters
		values = make(map[string]interface{}, 0)
		for idx := 0; idx < len(data)-1; idx += 2 {
			if _, ok := data[idx].(string); !ok {
				return 0, errors.New("unknow type for update")
			}
			values[data[idx].(string)] = data[idx+1]
		}
	}

	s.statement.Set(UPDATE, s.Schema().Name, values)
	// Build the complete statement
	sql, vars := s.statement.Build(UPDATE, WHERE)

	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Delete() (int64, error) {
	if s.Schema() == nil {
		return 0, errors.New("empty schema")
	}
	s.statement.Set(DELETE, s.Schema().Name)

	// Build the complete statement
	sql, vars := s.statement.Build(DELETE, WHERE)

	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Count() (int64, error) {
	if s.Schema() == nil {
		return 0, errors.New("empty schema")
	}
	s.statement.Set(COUNT, s.Schema().Name)
	// Build the complete statement
	sql, vars := s.statement.Build(COUNT, WHERE)

	row := s.Raw(sql, vars...).QueryRow()

	var cnt int64
	if err := row.Scan(&cnt); err != nil {
		return 0, err
	}

	return cnt, nil
}

func (s *Session) Limit(num int) *Session {
	s.statement.Set(LIMIT, num)
	return s
}

func (s *Session) Where(condition string, values ...interface{}) *Session {
	args := []interface{}{condition}
	args = append(args, values...)
	s.statement.Set(WHERE, args...)
	return s
}

func (s *Session) OrderBy(order string) *Session {
	s.statement.Set(ORDERBY, order)
	return s
}

func (s *Session) First(instance interface{}) error {
	elem := reflect.Indirect(reflect.ValueOf(instance))
	elemArr := reflect.New(reflect.SliceOf(elem.Type())).Elem()

	if err := s.Limit(1).List(elemArr.Addr().Interface()); err != nil {
		return err
	}

	if elemArr.Len() == 0 {
		return errors.New("not found")
	}
	elem.Set(elemArr.Index(0))
	return nil
}
