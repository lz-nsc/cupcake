package schema

import (
	"fmt"
	"go/ast"
	"reflect"

	"github.com/lz-nsc/cupcake/orm/translator"
)

// Schema represonts a column of a table
type Field struct {
	Name string
	Type string
	Tag  string
}

// Schema represonts a table in the database
type Schema struct {
	Model  interface{}
	Name   string
	Fields map[string]*Field
}

func (s Schema) GetField(name string) *Field {
	return s.Fields[name]
}

// Parse turn a struct into schema
func Parse(record interface{}, trans translator.Translator) (*Schema, error) {
	modelType := reflect.Indirect(reflect.ValueOf(record)).Type()
	schema := &Schema{
		Model:  record,
		Name:   modelType.Name(),
		Fields: make(map[string]*Field),
	}
	if modelType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expect struct but got a %s", modelType.Kind())
	}
	// Iterate struct's field
	for i := 0; i < modelType.NumField(); i++ {
		f := modelType.Field(i)
		if !f.Anonymous && ast.IsExported(f.Name) {
			field := &Field{
				Name: f.Name,
				// Translate p's type to type in target database
				Type: trans.DataTypeOf(reflect.Indirect(reflect.New(f.Type))),
			}
			if val, ok := f.Tag.Lookup("cupcakeorm"); ok {
				field.Tag = val
			}

			schema.Fields[field.Name] = field
		}
	}
	return schema, nil
}
