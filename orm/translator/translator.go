package translator

import "reflect"

var translatorMap = map[string]Translator{}

// Each database can implement its own translater to translate
// golang type to type in the database
type Translator interface {
	DataTypeOf(typ reflect.Value) string
	// Queries for databases to check whether a table exists are different,
	// so translator is needed for table exists query
	TableExistSQL(name string) (string, []interface{})
}

// RegisterTranslator is for databases to register their translators
func RegisterTranslator(driverName string, trans Translator) {
	translatorMap[driverName] = trans
}

func GetTranslator(driver string) (trans Translator, ok bool) {
	trans, ok = translatorMap[driver]
	return
}
