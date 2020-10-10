package jsref

import (
	"fmt"
	"reflect"
	"syscall/js"
)

// Marshal to js.Value
func Marshal(i interface{}) (js.Value, error) {
	v := reflect.ValueOf(i)
	return marshal(v)
}

func marshal(v reflect.Value) (js.Value, error) {
	t := v.Type()
	switch {
	case isScalar(t):
		return marshalScalar(v)
	case isArray(t):
		return marshalArray(v)
	case isMap(t):
		return marshalMap(v)
	case isStruct(t):
		return marshalStruct(v)
	default:
		return js.Null(), fmt.Errorf("unknown type: %s", t.String())
	}
}

func marshalStruct(v reflect.Value) (js.Value, error) {
	t := v.Type()
	if isPtr(t) {
		if !v.IsValid() {
			return js.Null(), nil
		}
		t = t.Elem()
		v = v.Elem()
	}
	m := map[string]interface{}{}
	for i := 0; i < t.NumField(); i++ {
		fld := t.Field(i)
		fv := v.Field(i)
		ft := fld.Type
		name := fld.Name
		tag := parseTag(fld.Tag)
		if tag.ignore {
			continue
		}
		if tag.name != "" {
			name = tag.name
		}
		if isPtr(ft) {
			if fv.IsNil() {
				continue
			}
			fv = fv.Elem()
		}
		val, err := marshal(fv)
		if err != nil {
			return js.Null(), err
		}
		m[name] = val
	}

	return js.ValueOf(m), nil
}

func marshalMap(v reflect.Value) (js.Value, error) {
	t := v.Type()
	if isPtr(t) {
		if v.IsNil() {
			return js.Null(), nil
		}
		t = t.Elem()
		v = v.Elem()
	}
	kt := t.Key()
	if !isString(kt) {
		return js.Null(), fmt.Errorf("map key must be string not %s", kt)
	}
	if t.Elem().String() == "interface{}" {
		// already map[string]interface{} just return it
		return js.ValueOf(v.Interface()), nil
	}
	m := map[string]interface{}{}
	var err error
	for _, key := range v.MapKeys() {
		val := v.MapIndex(key)
		kn, err := marshalScalar(key)
		if err != nil {
			return js.Null(), err
		}
		name := kn.String()
		ret, err := marshal(val)
		if err != nil {
			return js.Null(), err
		}
		m[name] = ret
	}
	if err != nil {
		return js.Null(), err
	}
	return js.ValueOf(m), nil
}

func marshalArray(v reflect.Value) (js.Value, error) {
	t := v.Type()
	if isPtr(t) {
		if !v.IsValid() {
			return js.Null(), nil
		}
		v = v.Elem()
	}
	arr := []interface{}{}
	for i := 0; i < v.Len(); i++ {
		val, err := marshal(v.Index(i))
		if err != nil {
			return js.Null(), err
		}
		arr = append(arr, val)
	}
	return js.ValueOf(arr), nil
}

func marshalScalar(v reflect.Value) (js.Value, error) {
	t := v.Type()
	if isPtr(t) {
		if v.IsNil() {
			return js.Null(), nil
		}
		v = v.Elem()
	}
	return js.ValueOf(v.Interface()), nil
}
