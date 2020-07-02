package utils

import (
	"errors"
	"reflect"
)

//求并集
func Union(slice1, slice2 []string) []string {
	m := make(map[string]int)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 0 {
			slice1 = append(slice1, v)
		}
	}
	return slice1
}

//求交集
func Intersect(slice1, slice2 interface{}) (res []interface{}) {
	defer func() {
		if r := recover(); nil != r {
			res = make([]interface{}, 0)
		}
	}()

	s1Val := reflect.ValueOf(slice1)
	s2Val := reflect.ValueOf(slice2)

	if reflect.TypeOf(slice1).Kind() != reflect.Slice ||
		reflect.TypeOf(slice2).Kind() != reflect.Slice {
		panic(errors.New("<Intersect>: type error"))
	}

	m := make(map[interface{}]struct{})
	nn := make([]interface{}, 0)

	for i := 0; i < s1Val.Len(); i++ {
		m[s1Val.Index(i).Interface()] = struct{}{}
	}

	for i := 0; i < s2Val.Len(); i++ {
		if _, ok := m[s2Val.Index(i).Interface()]; ok {
			nn = append(nn, s2Val.Index(i).Interface())
		}
	}
	return nn
}
