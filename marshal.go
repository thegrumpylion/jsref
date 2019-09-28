package jsref

import (
	"fmt"
	"reflect"
	"syscall/js"
)

func Marshal(i interface{}) (js.Value, error) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	switch {
	case isScalar(t):
		return scalarToJSValue(v)
	case isArray(t):
		return arrayToJSValue(v)
	case isMap(t):
		return mapToJSValue(v)
	case isStruct(t):
		return structToJSValue(v)
	default:
		return js.Null(), fmt.Errorf("unknown type: %s", t.String())
	}
}

func structToJSValue(v reflect.Value) (js.Value, error) {
	t := v.Type()
	if isPtr(t) {
		if !v.IsValid() {
			return js.Null(), nil
		}
		t = t.Elem()
		v = v.Elem()
	}
	fmt.Println("struct", t)
	m := map[string]interface{}{}
	var err error
	for i := 0; i < t.NumField(); i++ {
		fld := t.Field(i)
		fv := v.Field(i)
		ft := fld.Type
		if isPtr(ft) {
			if !fv.IsValid() {
				continue
			}
			fv = fv.Elem()
			ft = ft.Elem()
		}
		switch {
		case isScalar(t):
			m[fld.Name], err = scalarToJSValue(v)
		case isArray(t):
			m[fld.Name], err = arrayToJSValue(v)
		case isMap(t):
			m[fld.Name], err = mapToJSValue(v)
		case isStruct(t):
			m[fld.Name], err = structToJSValue(v)
		default:
			return js.Null(), fmt.Errorf("unknown type: %s", ft.String())
		}
	}
	if err != nil {
		return js.Null(), err
	}
	return js.ValueOf(m), nil
}

func mapToJSValue(v reflect.Value) (js.Value, error) {
	t := v.Type()
	if isPtr(t) {
		if !v.IsValid() {
			return js.Null(), nil
		}
		t = t.Elem()
		v = v.Elem()
	}
	fmt.Println("map", t)
	kt := t.Elem()
	if !isString(kt) {
		return js.Null(), fmt.Errorf("map key must be string not %s", kt)
	}
	m := map[string]interface{}{}
	var err error
	for _, key := range v.MapKeys() {
		val := v.MapIndex(key)
		kn, err := scalarToJSValue(key)
		if err != nil {
			return js.Null(), err
		}
		kname := kn.String()
		switch {
		case isScalar(t):
			m[kname], err = scalarToJSValue(val)
		case isArray(t):
			m[kname], err = arrayToJSValue(val)
		case isMap(t):
			m[kname], err = mapToJSValue(val)
		case isStruct(t):
			m[kname], err = structToJSValue(val)
		default:
			return js.Null(), fmt.Errorf("unknown type: %s", t.String())
		}
		_ = kname
	}
	if err != nil {
		return js.Null(), err
	}
	return js.ValueOf(m), nil
}

func arrayToJSValue(v reflect.Value) (js.Value, error) {
	t := v.Type()
	if isPtr(t) {
		if !v.IsValid() {
			return js.Null(), nil
		}
		t = t.Elem()
		v = v.Elem()
	}
	fmt.Println("arr", t)
	arr := []interface{}{}
	var err error
	var val js.Value
	for i := 0; i < v.Len(); i++ {
		switch {
		case isScalar(t):
			val, err = scalarToJSValue(v)
		case isArray(t):
			val, err = arrayToJSValue(v)
		case isMap(t):
			val, err = mapToJSValue(v)
		case isStruct(t):
			val, err = structToJSValue(v)
		}
		if err != nil {
			return js.Null(), err
		}
		arr = append(arr, val)
	}
	return js.ValueOf(arr), nil
}

func scalarToJSValue(v reflect.Value) (js.Value, error) {
	t := v.Type()
	if isPtr(t) {
		t = t.Elem()
		v = v.Elem()
	}
	fmt.Println("scalar", t)
	return js.ValueOf(v.Interface()), nil
}
