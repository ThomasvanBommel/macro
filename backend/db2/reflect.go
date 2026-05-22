package db2

import "reflect"

// ToSlice converts a struct into a slice of its field values.
func ToSlice(in any) []any {
	v := reflect.ValueOf(in)

	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil
	}

	out := make([]any, 0, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		if field.CanInterface() {
			out = append(out, field.Interface())
		}
	}

	return out
}

// ToPtrSlice converts a struct into a slice of pointers to its field values.
func ToPtrSlice(in any) []any {
	v := reflect.ValueOf(in)

	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil
	}

	out := make([]any, 0, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		if field.CanAddr() && field.CanInterface() {
			out = append(out, field.Addr().Interface())
		}
	}

	return out
}
