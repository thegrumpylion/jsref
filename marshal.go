package jsref

import (
	"fmt"
	"reflect"
	"syscall/js"
)

// Marshal to js.Value
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
		case isScalar(ft):
			m[fld.Name], err = scalarToJSValue(fv)
		case isArray(ft):
			m[fld.Name], err = arrayToJSValue(fv)
		case isMap(ft):
			m[fld.Name], err = mapToJSValue(fv)
		case isStruct(ft):
			m[fld.Name], err = structToJSValue(fv)
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
		et := t.Elem()
		switch {
		case isScalar(et):
			m[kname], err = scalarToJSValue(val)
		case isArray(et):
			m[kname], err = arrayToJSValue(val)
		case isMap(et):
			m[kname], err = mapToJSValue(val)
		case isStruct(et):
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
	arr := []interface{}{}
	var err error
	var val js.Value
	et := t.Elem()
	for i := 0; i < v.Len(); i++ {
		switch {
		case isScalar(et):
			val, err = scalarToJSValue(v.Index(i))
		case isArray(et):
			val, err = arrayToJSValue(v.Index(i))
		case isMap(et):
			val, err = mapToJSValue(v.Index(i))
		case isStruct(et):
			val, err = structToJSValue(v.Index(i))
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
	return js.ValueOf(v.Interface()), nil
}
