package server

import (
	"reflect"

	"github.com/wolfeidau/exitus/pkg/api"
)

const (
	defaultQuery    = ""
	defaultPageSize = 100
	defaultOffset   = 0
)

func listArgs(q *api.Q, limit *api.Limit, offset *api.Offset) (string, int, int) {
	return toString(q, defaultQuery), toInt(limit, defaultPageSize), toInt(offset, defaultOffset)
}

func toString(v interface{}, defaultValue string) string {
	drv, _, _ := derefPointersZero(reflect.ValueOf(v))
	if drv.Kind() == reflect.String {
		return drv.String()
	}

	return defaultValue
}

func toInt(v interface{}, defaultValue int) int {
	drv, _, _ := derefPointersZero(reflect.ValueOf(v))

	var res int

	switch drv.Kind() {
	case reflect.Int64:
		res = int(drv.Int())
	case reflect.Int32:
		res = int(drv.Int())
	case reflect.Int:
		res = int(drv.Int())
	}

	if res == 0 {
		return defaultValue
	}
	return res
}

func derefPointersZero(rv reflect.Value) (drv reflect.Value, isPtr bool, isNilPtr bool) {
	for rv.Kind() == reflect.Ptr {
		isPtr = true
		if rv.IsNil() {
			isNilPtr = true
			rt := rv.Type().Elem()
			for rt.Kind() == reflect.Ptr {
				rt = rt.Elem()
			}
			drv = reflect.New(rt).Elem()
			return
		}
		rv = rv.Elem()
	}
	drv = rv
	return
}
