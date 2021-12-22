package diff

import (
	"errors"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/util/gconv"
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
	field := reflect.ValueOf(src)
	switch reflect.TypeOf(src).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < field.Len(); i++ {
			item := field.Index(i)
			switch item.Kind() {
			case reflect.Struct:
				if len(keys) == 0 {
					return false, errors.New("需要传结构体中的标签名")
				}
				fieldValue := item.FieldByName(keys[0])
				if !fieldValue.CanInterface() {
					return false, errors.New("转换成 interface 失败")
				}
				if fieldValue.Kind() == reflect.Ptr {
					fieldValue = fieldValue.Elem()
				}
				if reflect.TypeOf(target).Kind() != fieldValue.Kind() {
					return false, errors.New("目标值与查询值类型不一致,请先转换(例如 int 与 int64")
				}
				if fieldValue.Interface() == target {
					return true, nil
				}
			default:
				if item.Kind() == reflect.Ptr {
					item = item.Elem()
				}
				if !item.CanInterface() {
					return false, errors.New("转换成 interface 失败")
				}
				if reflect.TypeOf(target).Kind() != item.Kind() {
					return false, errors.New("目标值与查询值类型不一致,请先转换(例如 int 与 int64")
				}
				if item.Interface() == target {
					// 普通的 slice
					return true, nil
				}
			}
		}
		return false, nil
	case reflect.Map:
		return field.MapIndex(reflect.ValueOf(target)).Interface() != nil, nil
	case reflect.String:
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
	return false, errors.New("不支持的类型")
}

type InterfaceResult struct {
	src       interface{}
	target    interface{}
	targetKey string
	srcKey    string
	added     []interface{}
	updated   []interface{}
	deleted   []interface{}
}

//
// NewDiff
//  @Description: 接收到返回值后,需要自己断言一下,例如 xxx.([]int)
//  @param src 原始值,通常传数据库那一份即可
//  @param target 新的值,通常传 request 传进来的即可
//  @param keys 唯一判断 keys, 通常为 ID
//  @return interface{} added
//  @return interface{} update
//  @return interface{} deleted
//  @return error
//
func NewDiff(src, target interface{}, srcKey, targetKey string) (*InterfaceResult, error) {
	srcType := reflect.TypeOf(src)
	targetType := reflect.TypeOf(target)
	srcKind := srcType.Kind()
	targetKind := targetType.Kind()
	if srcKind != reflect.Array && srcKind != reflect.Slice {
		return nil, errors.New("diff 的类型需要为 slice 或 array,若只是想 diff 下是否被修改,使用 diff.IsEditAny 方法")
	}
	if targetKind != reflect.Array && targetKind != reflect.Slice {
		return nil, errors.New("diff 的类型需要为 slice 或 array,若只是想 diff 下是否被修改,使用 diff.IsEditAny 方法")
	}
	if len(srcKey) == 0 || len(targetKey) == 0 {
		return nil, errors.New("keys 必须指定")
	}
	result := InterfaceResult{
		src:       src,
		target:    target,
		srcKey:    srcKey,
		targetKey: targetKey,
	}
	if err := result.diff(); err != nil {
		return nil, err
	}
	return &result, nil
}

func (i *InterfaceResult) diff() error {
	srcValue := reflect.ValueOf(i.src)
	targetValue := reflect.ValueOf(i.target)
	var (
		// id 统计 map
		idDiffMap = map[interface{}]int{}
		// id - 对应实体的 map
		dataMap = map[interface{}]interface{}{}
		// len
		srcLength    = srcValue.Len()
		targetLength = targetValue.Len()
	)
	for j := 0; j < Max(srcLength, targetLength); j++ {
		if srcLength > j {
			value := srcValue.Index(j)
			if value.Kind() == reflect.Ptr {
				value = value.Elem()
			}
			if !value.CanInterface() {
				return errors.New("不能转换成 interface 的类型")
			}
			id := gconv.Int64(value.FieldByName(i.srcKey).Interface())
			idDiffMap[id]++
			dataMap[id] = value.Interface()
		}
		if targetLength > j {
			value := targetValue.Index(j)
			if value.Kind() == reflect.Ptr {
				value = value.Elem()
			}
			if !value.CanInterface() || !value.FieldByName(i.targetKey).CanInterface() {
				return errors.New("不能转换成 interface 的类型")
			}
			id := gconv.Int64(value.FieldByName(i.targetKey).Interface())
			if id == 0 {
				i.added = append(i.added, value.Interface())
				continue
			}
			idDiffMap[id]++
			dataMap[id] = value.Interface()
		}
	}
	for k, v := range idDiffMap {
		switch v {
		case 1:
			i.deleted = append(i.deleted, dataMap[k])
		case 2:
			i.updated = append(i.updated, dataMap[k])
		default:
			return errors.New("只支持最多出现两次的比较")
		}
	}
	return nil
}

//
// uniformType
//  @Description: 主要为了兼容项目各种不同的 int 类型,统一转换为 int64 并返回
//  @param u
//  @return interface{}
//
func uniformType(u interface{}) int64 {
	switch value := u.(type) {
	case int:
		return int64(value)
	case int8:
		return int64(value)
	case int16:
		return int64(value)
	case int32:
		return int64(value)
	case uint:
		return int64(value)
	case uint8:
		return int64(value)
	case uint16:
		return int64(value)
	case uint32:
		return int64(value)
	case uint64:
		return int64(value)
	case float32:
		return int64(value)
	case float64:
		return int64(value)
	}
	return gconv.Int64(u)
}

func (i *InterfaceResult) GetAdded() []interface{} {
	if i == nil {
		return nil
	}
	return i.added
}

func (i *InterfaceResult) GetUpdated() []interface{} {
	if i == nil {
		return nil
	}
	return i.updated
}

func (i *InterfaceResult) GetDeleted() []interface{} {
	if i == nil {
		return nil
	}
	return i.deleted
}

// IsEditAny 仅接受两个结构体互相校验,rule会默认匹配标签名,并自动匹配下划线/小驼峰/大驼峰;不支持未导出的 field,硬拿会报错
func IsEditAny(src interface{}, target interface{}, rule []string) (string, bool, error) {
	srcValue := reflect.ValueOf(src)
	targetValue := reflect.ValueOf(target)
	// 指针判断
	if srcValue.Kind() == reflect.Ptr {
		srcValue = srcValue.Elem()
	}
	if targetValue.Kind() == reflect.Ptr {
		targetValue = targetValue.Elem()
	}
	// 判断必须为结构体,否则下方会 panic
	if srcValue.Kind() != reflect.Struct || targetValue.Kind() != reflect.Struct {
		return "", false, errors.New("结构不一致或不为结构体")
	}
	for _, key := range rule {
		tempSrcValue := srcValue.FieldByName(key)
		if !tempSrcValue.IsValid() {
			tempSrcValue = srcValue.FieldByName(underlineToHump(key))
		}
		tempTargetValue := targetValue.FieldByName(key)
		if !tempTargetValue.IsValid() {
			tempTargetValue = targetValue.FieldByName(underlineToHump(key))
		}
		// 如果值是指针,则转换一下
		if tempSrcValue.Kind() == reflect.Ptr {
			tempSrcValue = tempSrcValue.Elem()
		}
		if tempTargetValue.Kind() == reflect.Ptr {
			tempTargetValue = tempTargetValue.Elem()
		}
		// 不能转换成 interface 的直接返回
		if !tempSrcValue.CanInterface() || !tempSrcValue.CanInterface() {
			return "", false, errors.New("结构体无法转换到 interface")
		}
		tempSrcInterface := tempSrcValue.Interface()
		tempTargetInterface := tempTargetValue.Interface()
		switch tempSrcValue.Kind() {
		case
			reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64,
			reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64:
			if gconv.Int64(tempSrcInterface) != gconv.Int64(tempTargetInterface) {
				return key, true, nil
			}
		case
			reflect.String,
			reflect.Bool,
			reflect.UnsafePointer,
			reflect.Chan,
			reflect.Interface,
			reflect.Array:
			if tempTargetInterface != tempSrcInterface {
				// 尝试转成 int 试一下,如果还不相同就是不相同了
				if gconv.Int64(tempTargetInterface) != gconv.Int64(tempSrcInterface) {
					return key, true, nil
				}
			}
		default:
			// 不能比较的类型
			return "", false, errors.New("无法比较的类型")
		}
	}
	return "", false, nil
}

// 下划线 / 小驼峰 转大驼峰
func underlineToHump(a string) string {
	res, err := gregex.ReplaceStringFuncMatch(`_[a-z]{1}`, a, func(match []string) string {
		if len(match) != 0 && len(match[0]) != 0 {
			return strings.ToUpper(match[0][1:2])
		}
		return match[0]
	})
	if err != nil {
		return ""
	}
	res = strings.ToUpper(res[0:1]) + res[1:]
	return res
}

func IsInMapIntInt(a map[int]int, b int) bool {
	_, ok := a[b]
	return ok
}
