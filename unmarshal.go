package jsref

import (
	"fmt"
	"reflect"
	"syscall/js"
)

// Unmarshal from js.Value
func Unmarshal(i interface{}, val js.Value) error {
	v := reflect.ValueOf(i)
	return unmarshal(v, val)
}

func unmarshal(v reflect.Value, val js.Value) error {
	t := v.Type()
	switch {
	case isScalar(t):
		return unmarshalScalar(v, val)
	case isArray(t):
		return unmarshalArr(v, val)
	case isMap(t):
		return unmarshalMap(v, val)
	case isStruct(t):
		return unmarshalStruct(v, val)
	default:
		return fmt.Errorf("unknown type: %s", t.String())
	}
}

func unmarshalStruct(v reflect.Value, val js.Value) error {
	if !IsValid(val) {
		return nil
	}
	if !IsObject(val) {
		return fmt.Errorf("val must be object for struct")
	}
	t := v.Type()
	if isPtr(t) {
		if v.IsNil() {
			v.Set(reflect.New(t.Elem()))
		}
		v = v.Elem()
		t = t.Elem()
	}
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
		jsVal := val.Get(name)
		if !IsValid(jsVal) {
			continue
		}
		if isPtr(ft) {
			if fv.IsNil() {
				fv.Set(reflect.New(ft.Elem()))
			}
		}
		err := unmarshal(fv, jsVal)
		if err != nil {
			return err
		}
	}
	return nil
}

func unmarshalMap(v reflect.Value, val js.Value) error {
	if !IsValid(val) {
		return nil
	}
	if !IsObject(val) {
		return fmt.Errorf("val must be object for map")
	}
	t := v.Type()
	if isPtr(t) {
		if v.IsNil() {
			v.Set(reflect.New(t.Elem()))
		}
		v = v.Elem()
		t = t.Elem()
	}
	et := t.Elem()
	m := reflect.MakeMap(t)
	for _, key := range ObjectKeys(val) {
		goVal := reflect.New(et).Elem()
		jsVal := val.Get(key)
		err := unmarshal(goVal, jsVal)
		if err != nil {
			return err
		}
		m.SetMapIndex(reflect.ValueOf(key), goVal)
	}
	v.Set(m)
	return nil
}

func unmarshalArr(v reflect.Value, val js.Value) error {
	if !IsValid(val) {
		return nil
	}
	if !IsArray(val) {
		return fmt.Errorf("val must be array")
	}
	t := v.Type()
	if isPtr(t) {
		if v.IsNil() {
			v.Set(reflect.New(t.Elem()))
		}
		t = t.Elem()
		v = v.Elem()
	}
	slc := reflect.MakeSlice(t, 0, val.Length())
	te := t.Elem()
	for i := 0; i < val.Length(); i++ {
		jsVal := val.Index(i)
		goVal := reflect.New(te).Elem()
		err := unmarshal(goVal, jsVal)
		if err != nil {
			return err
		}
		slc = reflect.Append(slc, goVal)
	}
	v.Set(slc)
	return nil
}

func unmarshalScalar(v reflect.Value, val js.Value) error {
	if !IsValid(val) {
		return nil
	}
	t := v.Type()
	if isPtr(t) {
		if v.IsNil() {
			v.Set(reflect.New(t.Elem()))
		}
		v = v.Elem()
	}
	switch {
	case isBool(t):
		v.SetBool(val.Bool())
	case isString(t):
		v.SetString(val.String())
	case isInt(t):
		v.SetInt(int64(val.Int()))
	case isFloat(t):
		v.SetFloat(val.Float())
	default:
		return fmt.Errorf("unknown type: %s", t.String())
	}
	return nil
}
