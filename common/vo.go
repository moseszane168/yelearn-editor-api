package common

type DeleteVO struct {
	IDS []int64 `json:"ids" form:"id" binding:"min=1"`
} // @name DeleteVO

type IDVO struct {
	ID int64 `json:"id" form:"id"`
} // @name IDVO

type IDSVO struct {
	IDS []int64 `json:"ids" form:"ids"`
} // @name IDSVO

type ExcelParseVO struct {
	Msg     string      `json:"msg"`     // 错误提示,没有错误为空字符串
	Vos     interface{} `json:"vos"`     // 显示数据
	Success int         `json:"success"` // 成功条数
	Fault   int         `json:"fault"`   // 失败条数
	Header  []string    `json:"header"`  // 数据头部
} // @name ExcelParseVO
