/**
 * 备件库存
 */

package spare

import (
	"crf-mold/base"
	"crf-mold/common"
	"crf-mold/common/constant"
	"crf-mold/common/validate"
	"crf-mold/controller/dict"
	"crf-mold/controller/excel"
	"crf-mold/dao"
	"crf-mold/model"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// @Tags 备件库存
// @Summary 入库
// @Accept json
// @Produce json
// @Param Body body InboundSpareRequestVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /spare/inbound [post]
func InboundSpareRequest(c *gin.Context) {
	var vo InboundSpareRequestVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	var spareRequestInfo model.SpareRequest
	base.CopyProperties(&spareRequestInfo, vo)

	userId := c.GetHeader(constant.USERID)
	spareRequestInfo.InboundTime = time.Now().UnixNano()
	spareRequestInfo.CreatedBy = userId
	spareRequestInfo.UpdatedBy = userId

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	if err := tx.Table("spare_request").Create(&spareRequestInfo).Error; err != nil {
		panic(err)
	}

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 备件库存
// @Summary 出库
// @Accept json
// @Produce json
// @Param Body body OutboundSpareRequestVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /spare/outbound [post]
func OutboundSpareRequest(c *gin.Context) {
	var vo OutboundSpareRequestVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	var saveBatches []model.SpareRequest
	vos := GetInboundSpareRequest(vo.SpareCode)
	userId := c.GetHeader(constant.USERID)
	needCount := vo.Count

	for i := 0; i < len(vos); i++ {
		item := vos[i]
		c := item.Count

		var spareRequest model.SpareRequest
		spareRequest.Type = vo.Type
		spareRequest.Remark = vo.Remark
		spareRequest.SpareCode = vo.SpareCode
		spareRequest.Location = item.Location
		spareRequest.CreatedBy = userId
		spareRequest.UpdatedBy = userId
		spareRequest.InboundTime = item.InboundTime // 出库记录出的那条，使用入库时间界定

		if c > needCount { // 当前库存大于需要库存，只出需要的数量
			spareRequest.Count = -needCount
			saveBatches = append(saveBatches, spareRequest)
			needCount = 0
			break
		} else if c < needCount { // 当前库存小于需要库存，把当前库存全出了
			spareRequest.Count = -c
			needCount = needCount - c
		} else { // 当前库存等于需要库存，把当前库存全出了
			spareRequest.Count = -c
			saveBatches = append(saveBatches, spareRequest)
			needCount = 0
			break
		}
		saveBatches = append(saveBatches, spareRequest)
	}

	// 没库存，无法出库
	if needCount != 0 {
		panic(base.ResponseEnum[base.SPARE_REQUEST_NOT_ENOUGH])
	}

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	if err := tx.Table("spare_request").CreateInBatches(&saveBatches, len(saveBatches)).Error; err != nil {
		panic(err)
	}

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 备件库存
// @Summary 备件库存分页
// @Accept json
// @Produce json
// @Param query query PageSpareRequestVO true "PageSpareRequestVO"
// @Param AuthToken header string false "Token"
// @Success 200 {object} spare.PageSpareRequestOutVO
// @Router /spare/request/page [get]
func PageSpareRequest(c *gin.Context) {
	var vo PageSpareRequestVO
	if err := c.ShouldBindQuery(&vo); err != nil {
		panic(base.ParamsErrorN())
	}

	tx := dao.GetConn()

	rawSql := `
		select
			sr.spare_code,
			si.name,
			si.flavor,
			si.material,
			sr.location,
			sum(sr.count) as count
		from spare_request sr
		left join spare_info si on si.code = sr.spare_code and si.is_deleted = 'N'
		where sr.is_deleted = 'N' and (si.code like concat('%',?,'%') or si.name like concat('%',?,'%'))
		group by sr.spare_code,si.name,sr.location,si.flavor,si.material
		having count > 0
	`
	var pageVos []PageSpareRequestOutVO
	page := base.PageWithRawSQL(tx, &pageVos, vo.GetCurrentPage(), vo.GetSize(), rawSql, vo.CodeOrName, vo.CodeOrName)

	if len(pageVos) == 0 {
		page.List = []interface{}{}
	} else {
		page.List = pageVos
	}

	c.JSON(http.StatusOK, base.Success(page))
}

// @Tags 备件库存
// @Summary 备件入库记录查看
// @Accept json
// @Produce json
// @Param query query SpareRequestRecordVO true "SpareRequestRecordVO"
// @Param AuthToken header string false "Token"
// @Success 200 {object} spare.SpareRequestRecordOutVO
// @Router /spare/inbound [get]
func OneSpareRequestInbound(c *gin.Context) {
	var vo SpareRequestRecordVO
	if err := c.ShouldBindQuery(&vo); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	var spareRequestinfo []model.SpareRequest
	tx := dao.GetConn().Table("spare_request")
	if vo.GmtCreatedBegin != "" {
		tx.Where("gmt_created > ?", vo.GmtCreatedBegin)
	}
	if vo.GmtCreatedEnd != "" {
		tx.Where("gmt_created < ?", vo.GmtCreatedEnd)
	}

	tx.Where("spare_code = ? and is_deleted = 'N' and count > 0", vo.Code).Find(&spareRequestinfo)

	if len(spareRequestinfo) > 0 {
		var vos []SpareRequestRecordOutVO

		for _, v := range spareRequestinfo {
			var vo SpareRequestRecordOutVO
			base.CopyProperties(&vo, v)
			vos = append(vos, vo)
		}

		c.JSON(http.StatusOK, base.Success(vos))
	} else {
		c.JSON(http.StatusOK, base.SuccessN())
	}
}

// @Tags 备件库存
// @Summary 备件入库记录分页
// @Accept json
// @Produce json
// @Param query query SpareRequestRecordPageVO true "SpareRequestRecordPageVO"
// @Param AuthToken header string false "Token"
// @Success 200 {object} spare.SpareRequestRecordOutVO
// @Router /spare/inbound/page [get]
func PageOneSpareRequestInbound(c *gin.Context) {
	var vo SpareRequestRecordPageVO
	if err := c.ShouldBindQuery(&vo); err != nil {
		panic(base.ParamsErrorN())
	}

	var spareRequestinfo []model.SpareRequest
	tx := dao.GetConn().Table("spare_request")
	if vo.GmtCreatedBegin != "" {
		tx.Where("gmt_created > ?", vo.GmtCreatedBegin)
	}
	if vo.GmtCreatedEnd != "" {
		tx.Where("gmt_created < ?", vo.GmtCreatedEnd)
	}

	tx.Where("spare_code = ? and is_deleted = 'N' and count > 0", vo.Code).Order("gmt_created desc")

	page := base.Page(tx, &spareRequestinfo, vo.GetCurrentPage(), vo.GetSize())

	pageList := page.List.(*[]model.SpareRequest)
	if len(*pageList) > 0 {
		var vos []SpareRequestRecordOutVO
		for _, v := range *pageList {
			var vo SpareRequestRecordOutVO
			base.CopyProperties(&vo, v)
			vos = append(vos, vo)
		}

		page.List = vos

		c.JSON(http.StatusOK, base.Success(page))
	} else {
		page.List = []interface{}{}
		c.JSON(http.StatusOK, base.Success(page))
	}
}

// @Tags 备件库存
// @Summary 备件出库记录查看
// @Accept json
// @Produce json
// @Param query query SpareRequestRecordVO true "SpareRequestRecordVO"
// @Param AuthToken header string false "Token"
// @Success 200 {object} spare.SpareRequestRecordOutVO
// @Router /spare/outbound [get]
func OneSpareRequestOutbound(c *gin.Context) {
	var vo SpareRequestRecordVO
	if err := c.ShouldBindQuery(&vo); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	var spareRequestinfo []model.SpareRequest
	tx := dao.GetConn().Table("spare_request")
	if vo.GmtCreatedBegin != "" {
		tx.Where("gmt_created > ?", vo.GmtCreatedBegin)
	}
	if vo.GmtCreatedEnd != "" {
		tx.Where("gmt_created < ?", vo.GmtCreatedEnd)
	}

	tx.Where("spare_code = ? and is_deleted = 'N' and count < 0", vo.Code).Order("gmt_created desc").Find(&spareRequestinfo)

	if len(spareRequestinfo) > 0 {
		var vos []SpareRequestRecordOutVO

		for _, v := range spareRequestinfo {
			var vo SpareRequestRecordOutVO
			base.CopyProperties(&vo, v)
			vo.Count = -vo.Count
			vos = append(vos, vo)
		}

		c.JSON(http.StatusOK, base.Success(vos))
	} else {
		c.JSON(http.StatusOK, base.SuccessN())
	}
}

// @Tags 备件库存
// @Summary 备件出库记录分页
// @Accept json
// @Produce json
// @Param query query SpareRequestRecordPageVO true "SpareRequestRecordPageVO"
// @Param AuthToken header string false "Token"
// @Success 200 {object} spare.SpareRequestRecordOutVO
// @Router /spare/outbound/page [get]
func PageOneSpareRequestOutbound(c *gin.Context) {
	var vo SpareRequestRecordPageVO
	if err := c.ShouldBindQuery(&vo); err != nil {
		panic(base.ParamsErrorN())
	}

	var spareRequestinfo []model.SpareRequest
	tx := dao.GetConn().Table("spare_request")
	if vo.GmtCreatedBegin != "" {
		tx.Where("gmt_created > ?", vo.GmtCreatedBegin)
	}
	if vo.GmtCreatedEnd != "" {
		tx.Where("gmt_created < ?", vo.GmtCreatedEnd)
	}

	tx.Where("spare_code = ? and is_deleted = 'N' and count < 0", vo.Code).Order("gmt_created desc")

	page := base.Page(tx, &spareRequestinfo, vo.GetCurrentPage(), vo.GetSize())

	pageList := page.List.(*[]model.SpareRequest)
	if len(*pageList) > 0 {
		var vos []SpareRequestRecordOutVO
		for _, v := range *pageList {
			var vo SpareRequestRecordOutVO
			base.CopyProperties(&vo, v)
			vo.Count = -vo.Count
			vos = append(vos, vo)
		}

		page.List = vos

		c.JSON(http.StatusOK, base.Success(page))
	} else {
		page.List = []interface{}{}
		c.JSON(http.StatusOK, base.Success(page))
	}
}

// @Tags 备件库存
// @Summary 备件库存查看
// @Accept json
// @Produce json
// @Param code query string true "备件编码"
// @Param AuthToken header string false "Token"
// @Success 200 {object} spare.PageSpareRequestOutVO
// @Router /spare/request [get]
func OneSpareRequest(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		panic(base.ParamsErrorN())
	}

	vos := ViewSpareRequest(code)
	if len(vos) == 0 {
		c.JSON(http.StatusOK, base.Success([]interface{}{}))
	} else {
		c.JSON(http.StatusOK, base.Success(vos))
	}
}

// 查询指定备件编码的库存
func ViewSpareRequest(code string) []PageSpareRequestOutVO {
	var spareRequestinfo []model.SpareRequest
	dao.GetConn().Table("spare_request").Where("is_deleted = 'N' and spare_code = ?", code).Group("spare_code").Group("location").Select("spare_code,location,sum(`count`) as count").Order("count").Having("count > 0").Find(&spareRequestinfo)

	var vos []PageSpareRequestOutVO
	if len(spareRequestinfo) > 0 {
		// 缓存备件基本资料
		spareInfoCache := make(map[string]model.SpareInfo)
		for _, v := range spareRequestinfo {
			var vo PageSpareRequestOutVO
			base.CopyProperties(&vo, v)

			if v, ok := spareInfoCache[vo.SpareCode]; ok { // 缓存已存在
				vo.Name = v.Name
				vo.Flavor = v.Flavor
				vo.Material = v.Material
			} else { // 不存在就去数据库查
				var spareInfo model.SpareInfo
				// 没找到就不管了,也不加到返回值里面
				if err := dao.GetConn().Table("spare_info").Where("code = ?", vo.SpareCode).First(&spareInfo).Error; err != nil {
					// Donothing
				} else {
					spareInfoCache[vo.SpareCode] = spareInfo // 放进缓存
					vo.Name = spareInfo.Name
					vo.Flavor = spareInfo.Flavor
					vo.Material = spareInfo.Material
				}
			}

			vos = append(vos, vo)
		}
	}

	return vos
}

// 获取入库库存
func GetInboundSpareRequest(code string) []PageSpareRequestOutVO {
	var spareRequestinfo []model.SpareRequest
	dao.GetConn().Table("spare_request").Where("is_deleted = 'N' and spare_code = ?", code).Group("spare_code").Group("location").Group("inbound_time").Select("inbound_time,spare_code,location,sum(`count`) as count").Order("inbound_time asc").Having("count > 0").Find(&spareRequestinfo)

	var vos []PageSpareRequestOutVO
	if len(spareRequestinfo) > 0 {
		// 缓存备件基本资料
		for _, v := range spareRequestinfo {
			var vo PageSpareRequestOutVO
			base.CopyProperties(&vo, v)
			vos = append(vos, vo)
		}
	}

	return vos
}

// @Tags 备件库存
// @Summary 解析Excel文件内容
// @Accept json
// @Produce json
// @Param file formData file true "文件"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /spare/request/parse [post]
func ParseSpareRequestExcel(c *gin.Context) {
	header, err := c.FormFile("file")
	if err != nil {
		panic(base.ParamsErrorN())
	}

	// 生成随机excel文件名
	dst := excel.GenerateExcelFile("spareRequest_parse")

	// 保存上传的excel文件
	if err := c.SaveUploadedFile(header, dst); err != nil {
		panic(err)
	}

	// 删除文件
	defer func() {
		err = os.Remove(dst)
		if err != nil {
			panic(err)
		}
	}()

	// 解析上传的excel文件
	var sb strings.Builder
	vos := common.ParseExcel(dst, SpareRequestImportVO{}, &sb, nil)

	// 对每条数据进行校验
	success := 0
	fault := 0
	var unionMap = make(map[string](map[interface{}]bool))
	for i := 0; i < len(vos); i++ {
		vo := vos[i].(SpareRequestImportVO)
		// 字段校验
		msg := validate.Validate(validate.ExcelHandler, vo, i, unionMap)
		if msg != nil {
			sb.WriteString(msg.Error())
			fault++
		} else {
			success++
		}

		// 存在校验,备件code
		code := vo.SpareCode
		var count int64
		dao.GetConn().Table("spare_info").Where("code = ? and is_deleted = 'N'", code).Count(&count)
		if count == 0 {
			msg := fmt.Sprintf("第%d行:【备件编码】不存在.", i+1)
			sb.WriteString(msg)
			fault++
			success--
		}
	}

	// 封装excel头部的所有字段信息
	var head []string
	tags := base.GetStructTagsForFieldIndexKey(SpareRequestImportVO{}, "excel")
	for _, v := range tags {
		h := v["excel"]
		head = append(head, h)
	}

	c.JSON(http.StatusOK, base.Success(&common.ExcelParseVO{
		Success: success,
		Fault:   fault,
		Msg:     sb.String(),
		Vos:     vos,
		Header:  head,
	}))
}

// @Tags 备件库存
// @Summary 导入解析后的Excel内容
// @Accept json
// @Produce json
// @Param body body SpareRequestImportParsedVOS true "SpareRequestImportParsedVOS"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /spare/request/excel [post]
func ImportSpareRequest(c *gin.Context) {
	var vo SpareRequestImportParsedVOS
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	d := vo.Data
	// 字典类型转换
	dictMapArr := base.GetStructTagsForFieldIndexKey(SpareRequestImportVO{}, "dict")
	fields := base.GetStructFields(SpareRequestImportVO{})

	// 这里有点坑，使用range循环，对v进行修改实际上不会修改实际数组中的值,必须使用数组下标指向v进行修改
	for index, v := range d {
		// 循环判断每个字段是否需要转换字典值
		for i, fieldName := range fields {
			// 是否字段为字典取值
			groupCode := dictMapArr[i]["dict"]

			if groupCode == "" {
				continue
			}

			// 设置传入值为字典值
			dictValue := base.GetStructField(&v, fieldName)
			dictKey := dict.GetKeyByValue(groupCode, dictValue.(string), "")
			base.SetStructField(&d[index], fieldName, dictKey)
		}
	}

	// 转成entity对象
	vos := base.CopyPropertiesList(reflect.TypeOf(model.SpareRequest{}), d)
	datas := vos.([]interface{})
	userId := c.GetHeader(constant.USERID)
	var batchSave []model.SpareRequest

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)
	// 设置entity对象中的其它字段值
	for _, v := range datas {
		vv := v.(model.SpareRequest)
		vv.CreatedBy = userId
		vv.UpdatedBy = userId
		batchSave = append(batchSave, vv)
	}

	tx.Table("spare_request").CreateInBatches(batchSave, len(batchSave))

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 备件库存
// @Summary 导出全部
// @Accept json
// @Produce json
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /spare/request/excel [get]
func ExportSpareRequest(c *gin.Context) {
	var res []model.SpareRequest
	tx := dao.GetConn().Table("spare_request")

	tx.Where("is_deleted = 'N'").Group("spare_code").Group("location").Select("spare_code,location,sum(`count`) as count").Find(&res)

	spareInfoCache := make(map[string]model.SpareInfo, len(res))
	var vos []SpareRequestExportVO
	for _, v := range res {
		// 跳过库存为0的
		if v.Count <= 0 {
			continue
		}

		var vo SpareRequestExportVO
		base.CopyProperties(&vo, v)

		if v, ok := spareInfoCache[vo.SpareCode]; ok { // 缓存已存在
			vo.Name = v.Name
			vo.Flavor = v.Flavor
			vo.Material = v.Material
		} else { // 不存在就去数据库查
			var spareInfo model.SpareInfo
			// 没找到就不管了,也不加到返回值里面
			if err := dao.GetConn().Table("spare_info").Where("code = ?", vo.SpareCode).First(&spareInfo).Error; err != nil {
				continue
			} else {
				spareInfoCache[vo.SpareCode] = spareInfo // 放进缓存
				vo.Name = spareInfo.Name
				vo.Flavor = spareInfo.Flavor
				vo.Material = spareInfo.Material
			}
		}
		vos = append(vos, vo)
	}

	fileName := excel.GenerateExcelFile(excel.SPARE_REQUEST)
	info := common.ExportExcel(fileName, "", vos)

	c.JSON(http.StatusOK, base.Success(info))
}
