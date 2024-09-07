package dict

type CreateDictGroupVO struct {
	Name string `json:"name" binding:"required"` // 组名称，全局唯一
	Code string `json:"code" binding:"required"` // 组code，全局唯一
} // @name CreateDictGroupVO

type UpdateDictGroupVO struct {
	ID   int64  `json:"id" binding:"required"`   // 组ID
	Code string `json:"code" binding:"required"` // 组code，全局唯一
	Name string `json:"name" binding:"required"` // 新的组名称
} // @name UpdateDictGroupVO

type CreateDictPropertyVO struct {
	Key       string `json:"key" binding:"required"`       // 字典key,全局唯一
	ValueCn   string `json:"valueCn"`                      // 字典值，中文
	ValueEn   string `json:"valueEn"`                      // 字典值，英文
	Order     uint32 `json:"order"`                        // 排序
	GroupCode string `json:"groupCode" binding:"required"` // 字典组Code
} // @name CreateDictPropertyVO

type UpdateDictPropertyVO struct {
	ID        int64  `json:"id" binding:"required"`        // 字典ID
	Key       string `json:"key" binding:"required"`       // 字典key,全局唯一
	ValueCn   string `json:"valueCn"`                      // 字典值，中文
	ValueEn   string `json:"valueEn"`                      // 字典值，英文
	Order     uint32 `json:"order"`                        // 排序
	GroupCode string `json:"groupCode" binding:"required"` // 字典组Code
} // @name UpdateDictPropertyVO

type DictAllVO struct {
	GroupCode string `json:"groupCode"` // 字典组code
	ID        int64  `json:"id"`        // 字典ID
	Key       string `json:"key"`       // 字典key,全局唯一
	ValueCn   string `json:"valueCn"`   // 字典值，中文
	ValueEn   string `json:"valueEn"`   // 字典值，英文
} // @name DictAllVO
