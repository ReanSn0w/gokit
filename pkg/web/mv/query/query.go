package query

import (
	"errors"
	"net/url"
	"reflect"
	"strconv"
)

// Интерфейс валидации данных в структуре
type Validate interface {
	Validate() error
}

// Заполняет структуру на основе query-параметров.
func Decode(values url.Values, data interface{}) error {
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
		if paramValues, ok := values[paramName]; ok && len(paramValues) > 0 {
			if err := setFieldValue(field, paramValues); err != nil {
				return err
			}
		}
	}

	if data, ok := data.(Validate); ok {
		if err := data.Validate(); err != nil {
			return err
		}
	}

	if data, ok := data.(Validate); ok {
		if err := data.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func setFieldValue(field reflect.Value, values []string) error {
	value := values[0]

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
	case reflect.Slice:
		sliceType := field.Type().Elem().Kind()
		slice := reflect.MakeSlice(field.Type(), len(values), len(values))
		for i, paramValue := range values {
			switch sliceType {
			case reflect.String:
				slice.Index(i).SetString(paramValue)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				num, err := strconv.ParseInt(paramValue, 10, 64)
				if err != nil {
					return err
				}
				slice.Index(i).SetInt(num)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				num, err := strconv.ParseUint(paramValue, 10, 64)
				if err != nil {
					return err
				}
				slice.Index(i).SetUint(num)
			case reflect.Float32, reflect.Float64:
				num, err := strconv.ParseFloat(paramValue, 64)
				if err != nil {
					return err
				}
				slice.Index(i).SetFloat(num)
			case reflect.Bool:
				boolean, err := strconv.ParseBool(paramValue)
				if err != nil {
					return err
				}
				slice.Index(i).SetBool(boolean)
			default:
				return errors.New("неподдерживаемый тип элемента среза: " + sliceType.String())
			}
		}
		field.Set(slice)
	default:
		return errors.New("неподдерживаемый тип поля: " + field.Kind().String())
	}

	return nil
}
