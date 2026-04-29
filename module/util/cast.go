package util

import "reflect"

func CastMapToStruct[T any](target map[string]any) T {
	resultPtr := reflect.New(reflect.TypeFor[T]())

	elem := resultPtr.Elem()
	typeOfT := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		field := typeOfT.Field(i)

		if val, ok := target[field.Name]; ok {
			fieldVal := elem.Field(i)

			if fieldVal.CanSet() {
				inputVal := reflect.ValueOf(val)

				if inputVal.Type().AssignableTo(fieldVal.Type()) {
					fieldVal.Set(inputVal)
				} else if inputVal.Type().ConvertibleTo(fieldVal.Type()) {
					fieldVal.Set(inputVal.Convert(fieldVal.Type()))
				}
			}
		}
	}

	return elem.Interface().(T)
}
