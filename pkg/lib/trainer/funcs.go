package trainer

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// ObjectToString - печатает структуру в строку
func ObjectToString(obj any) string {
	return printObjectRecursive(obj, 1)
}

// MakeSchema - Создает JSON схему из структуры на golang
func MakeSchema(st any) []byte {
	schema := structToSchema(reflect.TypeOf(st))
	data, _ := json.Marshal(schema)
	return data
}

func printObjectRecursive(obj any, level int) string {
	val := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)

	// Обрабатываем указатель на структуру (в том числе nil)
	for typ.Kind() == reflect.Ptr {
		if val.IsNil() {
			return "null"
		}
		val = val.Elem()
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return fmt.Sprintf("%v", val.Interface())
	}

	var sb strings.Builder
	headerPrefix := strings.Repeat("#", level)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)
		if i > 0 {
			sb.WriteString("\n\n")
		}
		sb.WriteString(fmt.Sprintf("%s %s\n", headerPrefix, field.Name))
		// Рекурсивно вызываем для struct или *struct
		k := fieldVal.Kind()
		if k == reflect.Struct || (k == reflect.Ptr && !fieldVal.IsNil() && fieldVal.Elem().Kind() == reflect.Struct) {
			sb.WriteString("\n")
			sb.WriteString(printObjectRecursive(fieldVal.Interface(), level+1))
		} else if k == reflect.Ptr && fieldVal.IsNil() {
			sb.WriteString("null")
		} else {
			sb.WriteString(fmt.Sprintf("%v", fieldVal.Interface()))
		}
	}
	return sb.String()
}

func structToSchema(t reflect.Type) map[string]any {
	res := map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	}

	var required []string
	properties := res["properties"].(map[string]any)

	// Dereference pointer if needed
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		jsonTag := field.Tag.Get("json")
		jsonField := strings.Split(jsonTag, ",")[0]
		if jsonField == "" {
			jsonField = field.Name
		}
		// skip fields with "-" tag
		if jsonField == "-" {
			continue
		}

		enumTag := field.Tag.Get("enum")

		typ, requiredField := typeToSchema(field.Type)

		// Если есть enum-тег — вставляем enum со значениями
		if enumTag != "" {
			// разбиваем по запятой, убираем возможные пробелы после запятой
			enumValues := []string{}
			for _, v := range strings.Split(enumTag, ",") {
				enumValues = append(enumValues, strings.TrimSpace(v))
			}
			// преобразуем в map[string]any — нужна мутабельность
			typMap := typ.(map[string]any)
			typMap["enum"] = enumValues
			typ = typMap
		}

		properties[jsonField] = typ
		if requiredField {
			required = append(required, jsonField)
		}
	}

	if len(required) > 0 {
		res["required"] = required
	}
	return res
}

func typeToSchema(t reflect.Type) (any, bool) {
	// Dereference pointer if needed
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return map[string]any{"type": "integer"}, true
	case reflect.String:
		return map[string]any{"type": "string"}, true
	case reflect.Bool:
		return map[string]any{"type": "boolean"}, true
	case reflect.Float32, reflect.Float64:
		return map[string]any{"type": "number"}, true
	case reflect.Slice, reflect.Array:
		items, _ := typeToSchema(t.Elem())
		return map[string]any{"type": "array", "items": items}, true
	case reflect.Struct:
		return structToSchema(t), true
	default:
		return map[string]any{"type": "string"}, true
	}
}
