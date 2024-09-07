/**
 * 通过反射取结构体或者数组切片对象的属性或值等属性的相关常用方法
 */

package base

import (
	"reflect"
	"strconv"
	"time"
)

/**
 * 获取指定结构体对象的指定字段的指定tag值，字段值作为key
 */
func GetStructTagsForFieldNameKey(obj interface{}, tag ...string) (tags map[string]map[string]string) {
	tags = make(map[string]map[string]string)

	typ := reflect.TypeOf(obj)
	// 不是结构体类型
	if reflect.Struct != typ.Kind() || len(tag) == 0 {
		return
	}

	srcType := reflect.TypeOf(obj)
	propertyNums := srcType.NumField()

	for i := 0; i < propertyNums; i++ {
		property := srcType.Field(i)
		for _, v := range tag {
			tags[property.Name][v] = property.Tag.Get(v)
		}
	}

	return
}

/**
 * 获取指定结构体对象的指定字段的指定tag值，索引作为key
 */
func GetStructTagsForFieldIndexKey(obj interface{}, tag ...string) (tags []map[string]string) {
	typ := reflect.TypeOf(obj)
	// 不是结构体类型
	if reflect.Struct != typ.Kind() || len(tag) == 0 {
		return
	}

	srcType := reflect.TypeOf(obj)
	propertyNums := srcType.NumField()

	for i := 0; i < propertyNums; i++ {
		property := srcType.Field(i)
		t := make(map[string]string, len(tag))
		for _, v := range tag {
			t[v] = property.Tag.Get(v)
		}

		tags = append(tags, t)
	}

	return
}

/**
 * 通过反射获取指定结构体对象的所有字段名
 */
func GetStructFields(obj interface{}) (fields []string) {

	typ := reflect.TypeOf(obj)
	// 不是结构体类型
	if reflect.Struct != typ.Kind() {
		return
	}

	srcType := reflect.TypeOf(obj)
	propertyNums := srcType.NumField()

	for i := 0; i < propertyNums; i++ {
		property := srcType.Field(i)
		fields = append(fields, property.Name)
	}

	return
}

/*
 * 返回指定vo或者vo数组中对象的所有字段和指定的tag以及所有值
 */
func GetVoFieldTypesValues(vos interface{}, tag ...string) (fields []string, tags [](map[string]string), values [][]interface{}) {

	// 处理数组和切片类型vos
	switch reflect.TypeOf(vos).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(vos)
		for i := 0; i < s.Len(); i++ {
			var row []interface{}
			item := s.Index(i).Interface()
			srcType, srcValue := reflect.TypeOf(item), reflect.ValueOf(item)

			if i == 0 {
				propertyNums := srcType.NumField()
				for i := 0; i < propertyNums; i++ {
					property := srcType.Field(i)
					propertyValue := srcValue.FieldByName(property.Name)
					m := make(map[string]string)
					for _, v := range tag {
						m[v] = property.Tag.Get(v)
					}
					tags = append(tags, m)

					fields = append(fields, property.Name)
					row = append(row, propertyValue.Interface())
				}
			} else {
				propertyNums := srcType.NumField()
				for i := 0; i < propertyNums; i++ {
					property := srcType.Field(i)
					propertyValue := srcValue.FieldByName(property.Name)
					row = append(row, propertyValue.Interface())
				}
			}
			values = append(values, row)
		}
		return
	}

	// 处理单个属性
	srcType, srcValue := reflect.TypeOf(vos), reflect.ValueOf(vos)
	propertyNums := srcType.NumField()

	for i := 0; i < propertyNums; i++ {
		property := srcType.Field(i)
		propertyValue := srcValue.FieldByName(property.Name)
		m := make(map[string]string)
		for _, v := range tag {
			m[v] = property.Tag.Get(v)
		}
		tags = append(tags, m)
		fields = append(fields, property.Name)
		values = append(values, []interface{}{propertyValue.Interface()})
	}
	return
}

/**
 * 通过反射获取指定struct对象的字段，传入指定结构体变量的指针
 */
func GetStructField(voPtr interface{}, fieldName string) interface{} {
	f := reflect.ValueOf(voPtr).Elem().FieldByName(fieldName)
	return f.Interface()
}

/**
 * 通过反射设置指定struct对象的字段，传入指定结构体变量的指针
 */
func SetStructField(voPtr interface{}, fieldName string, value string) error {
	f := reflect.ValueOf(voPtr).Elem().FieldByName(fieldName)
	inf := f.Interface()

	// 指针类型
	typ := reflect.ValueOf(inf)
	if value == "" && typ.Kind() == reflect.Ptr && typ.IsNil() {
		return nil
	}

	if f.CanSet() {
		// 将string转换为指定的字段类型
		switch inf.(type) {
		case float32:
			val, err := strconv.ParseFloat(value, 32)
			if err != nil {
				return err
			}
			f.Set(reflect.ValueOf(val))
		case float64:
			val, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}
			f.Set(reflect.ValueOf(val))
		case int64:
			val, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			f.Set(reflect.ValueOf(val))
		case int:
			val, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			f.Set(reflect.ValueOf(val))
		case string:
			f.Set(reflect.ValueOf(value))
		case Time:
			val, err := FormatBaseTime(value)
			if err != nil {
				// logrus.Panic("参数错误")
				return err
			}
			f.Set(reflect.ValueOf(val))
		case time.Time:
			val, err := FormatTime(value)
			if err != nil {
				return err
			}
			f.Set(reflect.ValueOf(val))
		}
	}

	return nil
}

/**
 * 通过反射获取切片中的每一个元素
 */
func GetSliceElements(s interface{}) (res []interface{}) {
	sliceValue := reflect.ValueOf(s)

	for i := 0; i < sliceValue.Len(); i++ {
		item := sliceValue.Index(i).Interface()
		res = append(res, item)
	}

	return
}
