package jsref

import (
	"reflect"
	"syscall/js"
	"unicode"
)

// IsValid checks for valid js Value
func IsValid(v js.Value) bool {
	return !(v.Type() == js.TypeUndefined || v.Type() == js.TypeNull)
}

// IsBool checks if js Value is bool
func IsBool(v js.Value) bool {
	return v.Type() == js.TypeBoolean
}

// IsString checks if js Value is string
func IsString(v js.Value) bool {
	return v.Type() == js.TypeString
}

// IsNumber checks if js Value is number
func IsNumber(v js.Value) bool {
	return v.Type() == js.TypeNumber
}

// IsArray checks if js Value is array
func IsArray(v js.Value) bool {
	return v.Type() == js.TypeObject &&
		js.Global().Get("Array").Call("isArray", v).Bool()
}

// IsObject checks if js Value is object
func IsObject(v js.Value) bool {
	return v.Type() == js.TypeObject &&
		!js.Global().Get("Array").Call("isArray", v).Bool()
}

// ObjectKeys returns the keys of an object
func ObjectKeys(v js.Value) []string {
	out := []string{}
	keys := js.Global().Get("Object").Call("keys", v)
	for i := 0; i < keys.Length(); i++ {
		out = append(out, keys.Index(i).String())
	}
	return out
}

// go type matching

func isPtr(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr
}

func isBool(t reflect.Type) bool {
	if isPtr(t) {
		t = t.Elem()
	}
	return t.Kind() == reflect.Bool
}

func isInt(t reflect.Type) bool {
	if isPtr(t) {
		t = t.Elem()
	}
	return t.Kind() == reflect.Int ||
		t.Kind() == reflect.Int8 ||
		t.Kind() == reflect.Int16 ||
		t.Kind() == reflect.Int32 ||
		t.Kind() == reflect.Int64 ||
		t.Kind() == reflect.Uint ||
		t.Kind() == reflect.Uint8 ||
		t.Kind() == reflect.Uint16 ||
		t.Kind() == reflect.Uint32 ||
		t.Kind() == reflect.Uint64
}

func isFloat(t reflect.Type) bool {
	if isPtr(t) {
		t = t.Elem()
	}
	return t.Kind() == reflect.Float32 ||
		t.Kind() == reflect.Float64
}

func isString(t reflect.Type) bool {
	if isPtr(t) {
		t = t.Elem()
	}
	return t.Kind() == reflect.String
}

func isStruct(t reflect.Type) bool {
	if isPtr(t) {
		t = t.Elem()
	}
	return t.Kind() == reflect.Struct
}

func isMap(t reflect.Type) bool {
	if isPtr(t) {
		t = t.Elem()
	}
	return t.Kind() == reflect.Map
}

func isArray(t reflect.Type) bool {
	if isPtr(t) {
		t = t.Elem()
	}
	return t.Kind() == reflect.Slice ||
		t.Kind() == reflect.Array
}

func isNumber(t reflect.Type) bool {
	return isInt(t) || isFloat(t)
}

func isScalar(t reflect.Type) bool {
	return isBool(t) ||
		isNumber(t) || isString(t)
}

// string funcs

func upFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func lowFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}
