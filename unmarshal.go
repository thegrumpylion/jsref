package jsref

import (
	"fmt"
	"reflect"
	"syscall/js"
)

// Unmarshal from js.Value
func Unmarshal(i interface{}, val js.Value) error {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
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
		jsVal := val.Get(fld.Name)
		if !IsValid(jsVal) {
			jsVal = val.Get(lowFirst(fld.Name))
		}
		if !IsValid(jsVal) {
			return nil
		}
		switch {
		case isScalar(ft):
			err = unmarshalScalar(fv, jsVal)
		case isArray(ft):
			err = unmarshalArr(fv, jsVal)
		case isMap(ft):
			err = unmarshalMap(fv, jsVal)
		case isStruct(ft):
			err = unmarshalStruct(fv, jsVal)
		default:
			return fmt.Errorf("unknown type: %s", t.String())
		}
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
	var err error

	for _, key := range ObjectKeys(val) {
		goVal := reflect.New(et).Elem()
		jsVal := val.Get(key)
		switch {
		case isScalar(et):
			err = unmarshalScalar(goVal, jsVal)
		case isArray(et):
			err = unmarshalArr(goVal, jsVal)
		case isMap(et):
			err = unmarshalMap(goVal, jsVal)
		case isStruct(et):
			err = unmarshalStruct(goVal, jsVal)
		default:
			return fmt.Errorf("unknown type: %s", t.String())
		}
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
	slc := reflect.MakeSlice(t, 0, 10)
	et := t.Elem()
	var err error

	for i := 0; i < val.Length(); i++ {
		jsVal := val.Index(i)
		goVal := reflect.New(et).Elem()
		switch {
		case isScalar(et):
			err = unmarshalScalar(goVal, jsVal)
		case isArray(et):
			err = unmarshalArr(goVal, jsVal)
		case isMap(et):
			err = unmarshalMap(goVal, jsVal)
		case isStruct(et):
			err = unmarshalStruct(goVal, jsVal)
		default:
			return fmt.Errorf("unknown type: %s", t.String())
		}
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
