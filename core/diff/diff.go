package diff

import (
	"errors"
	"reflect"
	"strings"
)

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func IsIn(src interface{}, target interface{}, keys ...string) (bool, error) {
	switch reflect.TypeOf(src).Kind() {
	case reflect.Slice, reflect.Array:
		return isInSlice(src, target, keys...)
	case reflect.Map:
		return isInMap(src, target, keys...)
	case reflect.String:
		return isInString(src, target)
	}
	return false, errors.New("不支持的类型")
}

func isInMap(src interface{}, target interface{}, keys ...string) (bool, error) {
	field := reflect.ValueOf(src)
	if field.IsZero() || field.IsNil() {
		return false, nil
	}
	return field.MapIndex(reflect.ValueOf(target)).Interface() != nil, nil
}

func isInSlice(src interface{}, target interface{}, keys ...string) (bool, error) {
	field := reflect.ValueOf(src)
	for i := 0; i < field.Len(); i++ {
		item := field.Index(i)
		switch item.Kind() {
		case reflect.Struct:
			if len(keys) == 0 {
				return false, errors.New("缺少结构体中的标签名")
			}
			if reflect.ValueOf(item.Interface()).FieldByName(keys[0]).Interface() == target {
				return true, nil
			}
		default:
			if item.Interface() == target {
				// 普通的 slice
				return true, nil
			}
		}
	}
	return false, nil
}

func isInString(src interface{}, target interface{}) (bool, error) {
	strSrc, ok := src.(string)
	if !ok {
		return false, errors.New("src断言失败")
	}
	strTarget, ok := target.(string)
	if !ok {
		return false, errors.New("src 为 string 类型的时候,target 也需要为 string")
	}
	return strings.Contains(strSrc, strTarget), nil
}
