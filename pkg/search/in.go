package search

import (
	"errors"
	"fmt"
	"reflect"
)

func In(haystack interface{}, needle interface{}) (bool, error) {

	sVal := reflect.ValueOf(haystack)
	kind := sVal.Kind()

	if kind == reflect.Slice || kind == reflect.Array {
		for i := 0; i < sVal.Len(); i++ {
			fmt.Println(sVal.Index(i).Interface())
			if sVal.Index(i).Interface() == needle {
				return true, nil
			}
		}

		return false, nil
	}

	return false, errors.New("ErrUnSupportHaystack")
}