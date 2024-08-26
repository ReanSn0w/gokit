package web

import (
	"errors"
	"net/url"
	"reflect"
	"strconv"
)

// DecodeQuery заполняет структуру на основе query-параметров.
func DecodeQuery(values url.Values, data interface{}) error {
	// Проверьте что data является указателем на структуру
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("data должен быть непустым указателем на структуру")
	}

	// Получите значение, на которое указывает data
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("data должен быть указателем на структуру")
	}

	// Проведите обход полей структуры
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Пропустите приватные поля
		if !field.CanSet() {
			continue
		}

		// Получите имя параметра из тега "query" или используйте имя поля
		paramName := fieldType.Tag.Get("query")
		if paramName == "" {
			paramName = fieldType.Name
		}

		// Получите значение параметра
		if values, ok := values[paramName]; ok && len(values) > 0 {
			value := values[0]

			// Заполните поле в зависимости от его типа
			switch field.Kind() {
			case reflect.String:
				field.SetString(value)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				num, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return err
				}
				field.SetInt(num)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				num, err := strconv.ParseUint(value, 10, 64)
				if err != nil {
					return err
				}
				field.SetUint(num)
			case reflect.Float32, reflect.Float64:
				num, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return err
				}
				field.SetFloat(num)
			case reflect.Bool:
				boolean, err := strconv.ParseBool(value)
				if err != nil {
					return err
				}
				field.SetBool(boolean)
			default:
				return errors.New("неподдерживаемый тип поля: " + field.Kind().String())
			}
		}
	}

	return nil
}
