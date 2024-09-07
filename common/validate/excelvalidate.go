package validate

import (
	"fmt"
	"strconv"
)

var ExcelHandler *ExcelValidate = &ExcelValidate{}

type ExcelValidate struct {
}

func (p *ExcelValidate) Required(arg *ValidateArguments) string {
	return fmt.Sprintf("第%d行:【%s】不能为空.", arg.Index+1, arg.ExcelTagName)
}

func (p *ExcelValidate) Unique(arg *ValidateArguments) string {
	return fmt.Sprintf("第%d行:【%s】不能不能重复.", arg.Index+1, arg.ExcelTagName)
}

func (p *ExcelValidate) Min(arg *ValidateArguments) string {
	v := arg.Value
	vi, err := strconv.Atoi(arg.ValidateValue)
	if err != nil {
		panic(err)
	}

	switch vv := v.(type) {
	case string:
		if len(vv) < vi {
			return fmt.Sprintf("第%d行:【%s】长度最短%d个字符.", arg.Index+1, arg.ExcelTagName, vi)
		}
	case int64:

		if vv < int64(vi) {
			return fmt.Sprintf("第%d行:【%s】最小值%d.", arg.Index+1, arg.ExcelTagName, vi)
		}
	case int:
		if vv < vi {
			return fmt.Sprintf("第%d行:【%s】最小值%d.", arg.Index+1, arg.ExcelTagName, vi)
		}
	}
	return ""
}

func (p *ExcelValidate) Max(arg *ValidateArguments) string {

	vi, err := strconv.Atoi(arg.ValidateValue)
	if err != nil {
		panic(err)
	}

	switch vv := arg.Value.(type) {
	case string:
		if len(vv) > vi {
			return fmt.Sprintf("第%d行:【%s】长度最长%d个字符.", arg.Index+1, arg.ExcelTagName, vi)
		}
	case int64:
		if vv > int64(vi) {
			return fmt.Sprintf("第%d行:【%s】最大值%d.", arg.Index+1, arg.ExcelTagName, vi)
		}
	case int:
		if vv > vi {
			return fmt.Sprintf("第%d行:【%s】最大值%d.", arg.Index+1, arg.ExcelTagName, vi)
		}
	}

	return ""
}

func (p *ExcelValidate) Size(arg *ValidateArguments) string {
	return ""
}

func (p *ExcelValidate) Enum(arg *ValidateArguments) string {
	return fmt.Sprintf("第%d行：【%s】错误,可选值只能为:%s", arg.Index+1, arg.ExcelTagName, arg.DictValues)
}
