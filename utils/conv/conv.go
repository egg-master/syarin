package conv

import "reflect"

func StructToMap(val interface{}) (mapVal map[string]interface{}, ok bool) {
	structVal := reflect.Indirect(reflect.ValueOf(val))
	typ := structVal.Type()

	mapVal = make(map[string]interface{})

	for i := 0; i < typ.NumField(); i++ {
		field := structVal.Field(i)

		if field.CanSet() {
			mapVal[typ.Field(i).Name] = field.Interface()
		}
	}

	return
}

func MapToStruct(mapVal map[string]interface{}, val interface{}) (ok bool) {
	structVal := reflect.Indirect(reflect.ValueOf(val))
	for name, elem := range mapVal {
		structVal.FieldByName(name).Set(reflect.ValueOf(elem))
	}

	return
}
