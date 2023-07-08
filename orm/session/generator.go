package session

import (
	"fmt"
	"strings"
)

// Generator is used to construct a clause based on given action
type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[SELECT] = _select
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy
}

// TODO: Placeholder for mysql, sql is "?" while in PostgreSQL it is "$N"
// 		 We are using "?" here first, need to modify when PostgreSQL is integrated
func genHolderList(count int) (str string) {
	holders := make([]string, count)
	for idx := 0; idx < len(holders); idx++ {
		holders[idx] = "?"
	}
	return strings.Join(holders, ", ")
}

// SELECT column1, column2, ...
// FROM table_name;
func _select(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %v FROM %s", fields, tableName), []interface{}{}
}

// INSERT INTO table_name (column1, column2, column3, ...)
func _insert(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%s)", tableName, fields), []interface{}{}

}

//	VALUES
//		(value_list_1),
//		...
//		(value_list_n)
func _values(values ...interface{}) (string, []interface{}) {
	rows := []string{}
	var vars []interface{}

	for _, value := range values {
		valueList := value.([]interface{})
		holders := genHolderList(len(valueList))
		rows = append(rows, fmt.Sprintf("(%v)", holders))
		vars = append(vars, valueList...)
	}

	return fmt.Sprintf("VALUES %s", strings.Join(rows, ", ")), vars
}

func _where(values ...interface{}) (string, []interface{}) {
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

func _limit(values ...interface{}) (string, []interface{}) {
	return "LIMIT ?", values
}

func _orderBy(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}
