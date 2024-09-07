/**
 * 对象浅拷贝
 */

package base

import (
	"reflect"
	"time"
)

/**
 * 将源对象数据数组拷贝到新对象数组中,TODO:interface数组转指定类型
 */
func CopyPropertiesList(dstType reflect.Type, src interface{}) interface{} {

	var res []interface{}
	srcValue := reflect.ValueOf(src)
	// src支持结构体和指针类型
	if srcValue.Kind() != reflect.Slice && srcValue.Kind() == reflect.Ptr {
		srcValue = srcValue.Elem()
	}

	l := srcValue.Len()
	for i := 0; i < l; i++ {
		srcOriValue := srcValue.Index(i).Interface()
		dstPtr := reflect.New(dstType)
		dstValue := dstPtr.Interface()

		CopyProperties(dstValue, srcOriValue)
		res = append(res, dstPtr.Elem().Interface())
	}

	return res
}

/**
 * 不同结构之间拷贝属性
 */
func CopyProperties(dst, src interface{}) {
	dstType, dstValue := reflect.TypeOf(dst), reflect.ValueOf(dst)
	srcType, srcValue := reflect.TypeOf(src), reflect.ValueOf(src)

	// dst必须结构体指针类型
	if dstType.Kind() != reflect.Ptr || dstType.Elem().Kind() != reflect.Struct {
		panic("dst type should be a struct pointer")
	}

	// src必须为结构体或者结构体指针
	if srcType.Kind() == reflect.Ptr {
		srcType, srcValue = srcType.Elem(), srcValue.Elem()
	}
	if srcType.Kind() != reflect.Struct {
		panic("src type should be a struct or a struct pointer")
	}

	// 取具体内容
	dstType, dstValue = dstType.Elem(), dstValue.Elem()

	// 属性个数
	propertyNums := dstType.NumField()

	for i := 0; i < propertyNums; i++ {
		// 属性
		property := dstType.Field(i)
		// 待填充属性值
		propertyValue := srcValue.FieldByName(property.Name)

		// 无效
		if !propertyValue.IsValid() {
			continue
		}

		// 属性同名但类型不同
		if property.Type != propertyValue.Type() {
			// string类型转为Time
			if property.Type == reflect.TypeOf(time.Now()) && propertyValue.Type() == reflect.TypeOf("") {
				timeStr := propertyValue.String()
				time, err := FormatTime(timeStr)
				if err != nil {
					continue
				}
				if dstValue.Field(i).CanSet() {
					dstValue.Field(i).Set(reflect.ValueOf(time))
				}
			}
			// string类型转为base.Time
			if property.Type == reflect.TypeOf(Now()) && propertyValue.Type() == reflect.TypeOf("") {
				timeStr := propertyValue.String()
				time, err := FormatTime(timeStr)
				if err != nil {
					continue
				}
				if dstValue.Field(i).CanSet() {
					dstValue.Field(i).Set(reflect.ValueOf(time))
				}
			}

			// base.Time类型转为string
			if property.Type == reflect.TypeOf("") && propertyValue.Type() == reflect.TypeOf(Now()) {
				timeBase := propertyValue.Interface()
				t := timeBase.(Time)
				if dstValue.Field(i).CanSet() {
					dstValue.Field(i).Set(reflect.ValueOf(t.String()))
				}
			}

			// time.Time类型转为base.Time
			if property.Type == reflect.TypeOf(Now()) && propertyValue.Type() == reflect.TypeOf(time.Now()) {
				timeBase := propertyValue.Interface()
				t := timeBase.(time.Time)
				if dstValue.Field(i).CanSet() {
					dstValue.Field(i).Set(reflect.ValueOf(Time(t)))
				}
			}

			continue
		}

		if dstValue.Field(i).CanSet() {
			dstValue.Field(i).Set(propertyValue)
		}
	}
}
