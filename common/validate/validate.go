/**
 * 参数校验工具类
 */

package validate

import (
	"crf-mold/base"
	"crf-mold/dao"
	"crf-mold/model"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	REQUIRED = "required" // 必填
	UNIQUE   = "unique"   // 唯一
	MIN      = "min"      // 最小
	MAX      = "max"      // 最大
	SIZE     = "size"     // 长度
	ENUM     = "enum"     // 字典枚举key
	VALIDATE = "validate"
)

type ValidateArguments struct {
	Index         int         // 索引，如果是校验多个vo则传入，否则不传入
	DictValues    string      // 字典取值，使用","隔开，该值通过访问数据字典得到
	Value         interface{} // 通过反射取出的当前字段值
	ValidateValue string      // 校验值，如果是min=3，则校验值为3
	ExcelTagName  string      // 通过反射取出的excel tag 字段名
	FieldName     string      // 字段名称
}

type ValidateHandler interface {
	Required(arg *ValidateArguments) string
	Unique(arg *ValidateArguments) string
	Min(arg *ValidateArguments) string
	Max(arg *ValidateArguments) string
	Size(arg *ValidateArguments) string
	Enum(arg *ValidateArguments) string
}

type ValidateItem struct {
	Key   string
	Value string
} // @name ValidateItem

/**
 * 校验方法
 */
func Validate(handler ValidateHandler, vo interface{}, index int, unionMap map[string]map[interface{}]bool) error {
	val := reflect.ValueOf(vo)

	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		panic("参数错误")
	}

	var sb strings.Builder

	fields, tags, values := base.GetVoFieldTypesValues(vo, VALIDATE, "excel")

	for i := 0; i < len(fields); i++ {
		v := values[i][0]
		tag := tags[i]
		if tag == nil {
			continue
		}
		tagStr := tag[VALIDATE]
		excelTagStr := tag["excel"]
		if tagStr == "" {
			continue
		}
		if excelTagStr == "" {
			excelTagStr = fields[i]
		}

		validateItem := FormatValidateTag(tag[VALIDATE])
		for j := 0; j < len(validateItem); j++ {
			validateKey := validateItem[j].Key
			validateValue := validateItem[j].Value

			switch validateKey {
			case REQUIRED:
				if base.IsZeroValue(v) {
					m := (handler).Required(&ValidateArguments{
						Index:         index,
						DictValues:    "",
						Value:         v,
						ValidateValue: "",
						ExcelTagName:  excelTagStr,
						FieldName:     fields[i],
					})

					if m != "" {
						sb.WriteString(m)
						sb.WriteString("\n")
					}
				}
			case UNIQUE:
				if unionMap == nil {
					break
				}
				m := unionMap[fields[i]]
				if m == nil {
					m = make(map[interface{}]bool)
					unionMap[fields[i]] = m
				}

				if !unionMap[fields[i]][v] {
					unionMap[fields[i]][v] = true
				} else {
					m := (handler).Unique(&ValidateArguments{
						Index:         index,
						DictValues:    "",
						Value:         v,
						ValidateValue: "",
						ExcelTagName:  excelTagStr,
						FieldName:     fields[i],
					})

					if m != "" {
						sb.WriteString(m)
						sb.WriteString("\n")
					}
				}
			case MIN:
				m := (handler).Min(&ValidateArguments{
					Index:         index,
					DictValues:    "",
					Value:         v,
					ValidateValue: validateValue,
					ExcelTagName:  excelTagStr,
					FieldName:     fields[i],
				})
				if m != "" {
					sb.WriteString(m)
					sb.WriteString("\n")
				}

			case MAX:
				m := (handler).Max(&ValidateArguments{
					Index:         index,
					DictValues:    "",
					Value:         v,
					ValidateValue: validateValue,
					ExcelTagName:  excelTagStr,
					FieldName:     fields[i],
				})
				if m != "" {
					sb.WriteString(m)
					sb.WriteString("\n")
				}

			case SIZE:
				m := (handler).Size(&ValidateArguments{
					Index:         index,
					DictValues:    "",
					Value:         v,
					ValidateValue: validateValue,
					ExcelTagName:  excelTagStr,
					FieldName:     fields[i],
				})
				if m != "" {
					sb.WriteString(m)
					sb.WriteString("\n")
				}
			case ENUM:
				// 字典key
				dictKey := validateValue
				var result []model.DictProperty
				dao.GetConn().Table("dict_property").Where("is_deleted = 'N' and group_code = ?", dictKey).Find(&result)
				flag := false
				optionStr := ""
				for _, item := range result {
					valueCn := item.ValueCn
					optionStr = optionStr + valueCn + ","
					if v == item.ValueCn {
						flag = true
						break
					}
				}

				if !flag {
					m := (handler).Enum(&ValidateArguments{
						Index:         index,
						DictValues:    optionStr[:len(optionStr)-1],
						Value:         v,
						ValidateValue: validateValue,
						ExcelTagName:  excelTagStr,
						FieldName:     fields[i],
					})

					if m != "" {
						sb.WriteString(m)
						sb.WriteString("\n")
					}
				}
			default:
				panic("validate type error:" + validateKey)
			}
		}
	}

	msg := sb.String()

	if msg == "" {
		return nil
	}
	return errors.New(msg)
}

/**
 * 校验指定入参VO对象,待删除
 */
func ValidateExcelVO(vo interface{}, index int, unionMap map[string]map[interface{}]bool) error {
	val := reflect.ValueOf(vo)

	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		panic("参数错误")
	}

	var sb strings.Builder

	fields, tags, values := base.GetVoFieldTypesValues(vo, VALIDATE, "excel")

	for i := 0; i < len(fields); i++ {
		v := values[i][0]
		tag := tags[i]
		if tag == nil {
			continue
		}
		tagStr := tag[VALIDATE]
		excelTagStr := tag["excel"]
		if tagStr == "" {
			continue
		}
		if excelTagStr == "" {
			excelTagStr = fields[i]
		}

		validateItem := FormatValidateTag(tag[VALIDATE])
		for j := 0; j < len(validateItem); j++ {
			validateKey := validateItem[j].Key
			validateValue := validateItem[j].Value

			switch validateKey {
			case REQUIRED:
				if base.IsZeroValue(v) {
					m := fmt.Sprintf("第%d行:【%s】字段不能为空.", index+1, excelTagStr)
					sb.WriteString(m)
					sb.WriteString("\n")
				}
			case UNIQUE:
				if unionMap == nil {
					break
				}
				m := unionMap[fields[i]]
				if m == nil {
					m = make(map[interface{}]bool)
					unionMap[fields[i]] = m
				}

				if !unionMap[fields[i]][v] {
					unionMap[fields[i]][v] = true
				} else {
					m := fmt.Sprintf("第%d行:【%s】字段不能不能重复.", index+1, excelTagStr)
					sb.WriteString(m)
					sb.WriteString("\n")
				}
			case MIN:
				vi, err := strconv.Atoi(validateValue)
				if err != nil {
					panic(err)
				}

				switch vv := v.(type) {
				case string:
					if len(vv) < vi {
						m := fmt.Sprintf("第%d行:【%s】字段长度最短%d个字符.", index+1, excelTagStr, vi)
						sb.WriteString(m)
						sb.WriteString("\n")
					}
				case int64:
					if vv < int64(vi) {
						m := fmt.Sprintf("第%d行:【%s】字段最小值%d.", index+1, excelTagStr, vi)
						sb.WriteString(m)
						sb.WriteString("\n")
					}
				case int:
					if vv < vi {
						m := fmt.Sprintf("第%d行:【%s】字段最小值%d.", index+1, excelTagStr, vi)
						sb.WriteString(m)
						sb.WriteString("\n")
					}
				}
			case MAX:
				vi, err := strconv.Atoi(validateValue)
				if err != nil {
					panic(err)
				}
				switch vv := v.(type) {
				case string:
					if len(vv) > vi {
						m := fmt.Sprintf("第%d行:【%s】字段长度最长%d个字符.", index+1, excelTagStr, vi)
						sb.WriteString(m)
						sb.WriteString("\n")
					}
				case int64:
					if vv > int64(vi) {
						m := fmt.Sprintf("第%d行:【%s】字段最大值%d.", index+1, excelTagStr, vi)
						sb.WriteString(m)
						sb.WriteString("\n")
					}
				case int:
					if vv > vi {
						m := fmt.Sprintf("第%d行:【%s】字段最大值%d.", index+1, excelTagStr, vi)
						sb.WriteString(m)
						sb.WriteString("\n")
					}
				}
			case SIZE:
			case ENUM:
				// 字典key
				dictKey := validateValue
				var result []model.DictProperty
				dao.GetConn().Table("dict_property").Where("is_deleted = 'N' and group_code = ?", dictKey).Find(&result)
				flag := false
				optionStr := ""
				for _, item := range result {
					valueCn := item.ValueCn
					optionStr = optionStr + valueCn + ","
					if v == item.ValueCn {
						flag = true
						break
					}
				}

				if !flag {
					m := fmt.Sprintf("第%d行：【%s】字段错误,可选值只能为:%s", index+1, excelTagStr, optionStr[:len(optionStr)-1])
					sb.WriteString(m)
					sb.WriteString("\n")
				}
			default:
				panic("validate type error")
			}
		}
	}

	msg := sb.String()

	if msg == "" {
		return nil
	}
	return errors.New(msg)
}

/**
 * 格式化返回结构体中自定义的validate tag数组
 */
func FormatValidateTag(tagStr string) (res []ValidateItem) {
	split := strings.Split(tagStr, ",")
	if len(split) == 0 {
		return
	}

	for _, v := range split {
		ones := strings.Split(v, "=")

		if len(ones) == 2 {
			res = append(res, ValidateItem{
				Key:   ones[0],
				Value: ones[1],
			})
		} else {
			res = append(res, ValidateItem{
				Key:   ones[0],
				Value: "",
			})
		}
	}

	return
}
