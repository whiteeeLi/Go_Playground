package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type mysql struct{}

var _ Dialect = (*mysql)(nil)

func init() {
	RegisterDialect("mysql", &mysql{})
}

func (m *mysql) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.Int, reflect.Int8, reflect.Int16:
		return "INT"
	case reflect.Int32:
		return "INTEGER"
	case reflect.Int64:
		return "BIGINT"
	case reflect.Uint, reflect.Uint8, reflect.Uint16:
		return "UNSIGNED INT"
	case reflect.Uint32:
		return "UNSIGNED INTEGER"
	case reflect.Uint64, reflect.Uintptr:
		return "UNSIGNED BIGINT"
	case reflect.Float32:
		return "FLOAT"
	case reflect.Float64:
		return "DOUBLE"
	case reflect.String:
		return "VARCHAR(255)" // 默认长度为 255，可根据需要调整
	case reflect.Array, reflect.Slice:
		return "BLOB"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "DATETIME"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

func (m *mysql) TableExistSQL(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?", args
}
