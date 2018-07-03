package maps // import "syreclabs.com/go/maps"

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

func PathSet(root interface{}, path string, value interface{}) error {
	if root == nil {
		return errors.New("nil obj")
	}
	typ := reflect.TypeOf(root)
	if typ.Kind() != reflect.Ptr {
		return errors.New("non-pointer obj: " + typ.String())
	}
	if typ.Elem().Kind() != reflect.Map {
		return errors.New("non-map obj: " + typ.String())
	}

	keys, err := parsePath(path)
	if err != nil {
		return err
	}

	// walk backwards through path and prepare container chain
	var val interface{} = value
	for i := len(keys); i > 0; i-- {
		k := keys[i-1]
		if idx, err := strconv.ParseInt(k, 10, 64); err == nil {
			// slice
			s := make([]interface{}, idx+1, idx+1)
			s[idx] = val
			val = s
		} else {
			// map
			val = map[string]interface{}{k: val}
		}
	}

	// merge prepared container chain with root
	merge(reflect.Indirect(reflect.ValueOf(root)).Interface(), val)

	return nil
}

var errInvalidPath = errors.New("invalid path")

func parsePath(path string) ([]string, error) {
	// support both "Foo[0].Bar" and "Foo.0.Bar" path syntax
	path = strings.Replace(path, "[", ".", -1)
	path = strings.Replace(path, "]", "", -1)

	keys := strings.Split(path, ".")
	if len(keys) == 0 {
		return nil, errInvalidPath
	}

	if keys[0] == "" || keys[0] == "." {
		return nil, errInvalidPath
	}
	if _, err := strconv.ParseInt(keys[0], 10, 64); err == nil {
		return nil, errInvalidPath
	}

	return keys, nil
}

func merge(dst, src interface{}) interface{} {
	switch s := src.(type) {
	case map[string]interface{}:
		var mdst map[string]interface{}
		mdst, ok := dst.(map[string]interface{})
		if !ok {
			mdst = make(map[string]interface{}, len(s))
		}
		for k := range s {
			mdst[k] = merge(mdst[k], s[k])
		}
		return mdst
	case []interface{}:
		var sdst []interface{}
		sdst, ok := dst.([]interface{})
		if !ok {
			sdst = make([]interface{}, len(s))
		}
		if len(sdst) < len(s) {
			ns := make([]interface{}, len(s))
			copy(ns, sdst)
			sdst = ns
		}
		for i := range s {
			if v := merge(sdst[i], s[i]); v != nil {
				sdst[i] = v
			}
		}
		return sdst
	default:
		return src
	}
}
