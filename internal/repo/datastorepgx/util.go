package datastorepgx

import (
	"reflect"

	"github.com/jackc/pgtype"
)

func textArrayToSlice(t pgtype.TextArray) []string {
	s := make([]string, 0, len(t.Elements))

	for _, v := range t.Elements {
		if v.Status == pgtype.Present {
			s = append(s, v.String)
		}
	}
	return s
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
