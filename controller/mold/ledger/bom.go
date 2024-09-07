/**
 * 模具文档
 */
package mold

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
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// @Tags 模具BOM
// @Summary 添加模具BOM
// @Accept json
// @Produce json
// @Param Body body CreateBomVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold/bom [post]
func CreateBom(c *gin.Context) {
	var vo CreateBomVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsErrorN())
	}

	// 参数校验,TODO

	// 模具BOM编码唯一
	var count int64
	dao.GetConn().Table("mold_bom").Where("standard_component = ? and part_code = ? and is_deleted = 'N'", vo.StandardComponent, vo.PartCode).Count(&count)
	if count > 0 {
		panic(base.ResponseEnum[base.MOLD_BOM_EXIST])
	}

	var moldBom model.MoldBom
	base.CopyProperties(&moldBom, vo)

	userId := c.GetHeader(constant.USERID)
	moldBom.CreatedBy = userId
	moldBom.UpdatedBy = userId

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	if err := tx.Table("mold_bom").Create(&moldBom).Error; err != nil {
		panic(err)
	}

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(moldBom.ID))
}

// @Tags 模具BOM
// @Summary 更新模具BOM
// @Accept json
// @Produce json
// @Param Body body UpdateBomVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold/bom [put]
func UpdateBom(c *gin.Context) {
	var vo UpdateBomVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsErrorN())
	}

	// 参数校验,TODO

	if vo.ID == 0 {
		panic(base.ParamsErrorN())
	}

	// 存在
	var one model.MoldBom
	if err := dao.GetConn().Table("mold_bom").Where("id = ?", vo.ID).First(&one).Error; err != nil {
		panic(base.ResponseEnum[base.MOLD_BOM_NOT_EXIST])
	}

	// 模具BOM编码唯一
	var count int64
	dao.GetConn().Table("mold_bom").Where("standard_component = ? and part_code = ? and is_deleted = 'N' and id != ?", one.StandardComponent, vo.PartCode, vo.ID).Count(&count)
	if count > 0 {
		panic(base.ResponseEnum[base.MOLD_BOM_EXIST])
	}

	// 更新
	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	var moldBom model.MoldBom
	base.CopyProperties(&moldBom, vo)
	moldBom.UpdatedBy = c.GetHeader(constant.USERID)

	if err := tx.Table("mold_bom").Where("id = ?", moldBom.ID).Updates(&moldBom).Error; err != nil {
		panic(err)
	}

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 模具BOM
// @Summary 用户模具BOM分页
// @Accept json
// @Produce json
// @Param query query PageBomVO true "PageBomVO"
// @Param AuthToken header string false "Token"
// @Success 200 {object} mold.PageBomOutVO
// @Router /mold/bom/page [get]
func PageBom(c *gin.Context) {
	var vo PageBomVO
	if err := c.ShouldBindQuery(&vo); err != nil {
		panic(base.ParamsErrorN())
	}

	var moldBom model.MoldBom
	base.CopyProperties(&moldBom, vo)

	tx := dao.GetConn().Table("mold_bom")
	if vo.CodeOrName == "" {
		tx = dao.BuildWhereCondition(tx, moldBom)
	} else {
		tx = tx.Where("(part_code like concat('%',?,'%') or part_name like concat('%',?,'%'))", vo.CodeOrName, vo.CodeOrName)
		tx = tx.Where("standard_component = ?", vo.StandardComponent)
		tx = tx.Where("mold_code = ?", vo.MoldCode)
	}

	tx = tx.Where("is_deleted = 'N'").Order("gmt_created desc")

	var result []model.MoldBom
	page := base.Page(tx, &result, vo.GetCurrentPage(), vo.GetSize())

	// 转换为outvo
	vos := page.List.(*[]model.MoldBom)
	pageVos, _ := base.CopyPropertiesList(reflect.TypeOf(PageBomOutVO{}), vos).([]interface{})
	var list []PageBomOutVO
	for _, v := range pageVos {
		vv, _ := v.(PageBomOutVO)
		list = append(list, vv)
	}

	if len(list) == 0 {
		page.List = []interface{}{}
	} else {
		page.List = list
	}

	c.JSON(http.StatusOK, base.Success(page))
}

// @Tags 模具BOM
// @Summary 用户模具BOM查看
// @Accept json
// @Produce json
// @Param id query string true "id"
// @Param AuthToken header string false "Token"
// @Success 200 {object} mold.PageBomOutVO
// @Router /mold/bom/one [get]
func OneBom(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		panic(base.ParamsErrorN())
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		panic(base.ParamsErrorN())
	}

	var moldinfo []model.MoldBom
	dao.GetConn().Table("mold_bom").Where("id = ? and is_deleted = 'N'", id).Find(&moldinfo)

	if len(moldinfo) > 0 {
		m := moldinfo[0]
		var vo PageBomOutVO
		base.CopyProperties(&vo, m)

		c.JSON(http.StatusOK, base.Success(vo))
	} else {
		c.JSON(http.StatusOK, base.SuccessN())
	}
}

// @Tags 模具BOM
// @Summary 删除模具BOM
// @Accept json
// @Produce json
// @Param Body body common.DeleteVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold/bom [delete]
func DeleteBom(c *gin.Context) {
	var vo common.DeleteVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsErrorN())
	}

	ids := vo.IDS
	if len(ids) == 0 {
		panic(base.ParamsErrorN())
	}

	// dao.BuildWhereIdIn(tx, ids)
	dao.GetConn().Table("mold_bom").Where("id in ?", ids).Updates(&model.MoldBom{
		IsDeleted: "Y",
		UpdatedBy: c.GetHeader(constant.USERID),
	})

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 模具BOM
// @Summary 解析Excel文件内容
// @Accept json
// @Produce json
// @Param file formData file true "文件"
// @Param standard formData string true "是否标准件,Y:标准件 N:非标准件"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold/bom/parse [post]
func ParseBomExcel(c *gin.Context) {
	header, err := c.FormFile("file")
	if err != nil {
		panic(base.ParamsErrorN())
	}
	standard := c.PostForm("standard")
	if standard == "" {
		panic(base.ParamsErrorN())
	}

	// 生成随机excel文件名
	dst := excel.GenerateExcelFile("bom_parse")

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

	// 对每条数据进行校验
	success := 0
	fault := 0

	var sb strings.Builder
	faultCountMap := make(map[int]string)
	// 解析上传的excel文件
	vos := common.ParseExcel(dst, BomImportVO{}, &sb, &faultCountMap)

	if len(faultCountMap) > 0 {
		fault = len(faultCountMap)
	}

	// 如果格式没有错误就进行校验
	if fault == 0 {
		var unionMap = make(map[string](map[interface{}]bool))
		for i := 0; i < len(vos); i++ {

			msg := validate.Validate(validate.ExcelHandler, vos[i], i, unionMap)
			if msg != nil {
				sb.WriteString(msg.Error())
				fault++
			} else {
				success++
			}
		}

		// 校验一下唯一性，上面的唯一性只是针对导入的数据进行的，这里的唯一性是针对表中的数据进行的
		// 上面如果有问题，这里就先不进行数据库操作校验了
		if fault == 0 {
			for i := 0; i < len(vos); i++ {
				errFlag := false
				// 类型强转
				v := vos[i].(BomImportVO)
				code := v.PartCode
				var count int64
				dao.GetConn().Table("mold_bom").Where("mold_code = ? and standard_component = ? and  part_code = ? and is_deleted = 'N'", v.MoldCode, standard, code).Count(&count)
				if count > 0 {
					m := fmt.Sprintf("第%d行:【零件号】字段已存在.", i+1)
					sb.WriteString(m)
					errFlag = true
				}
				// 类型强转
				dao.GetConn().Table("mold_info").Where("code = ? and is_deleted = 'N'", v.MoldCode).Count(&count)
				if count == 0 {
					m := fmt.Sprintf("第%d行:【模具编号】字段不存在.", i+1)
					sb.WriteString(m)
					errFlag = true
				}

				if errFlag {
					fault++
					success--
				}
			}
		}
	}

	// 封装excel头部的所有字段信息
	var head []string
	tags := base.GetStructTagsForFieldIndexKey(BomImportVO{}, "excel")
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

// @Tags 模具BOM
// @Summary 导入解析后的Excel内容
// @Accept json
// @Produce json
// @Param body body BomImportParsedVOS true "BomImportParsedVOS"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold/bom/excel [post]
func ImportBom(c *gin.Context) {
	var vo BomImportParsedVOS
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsErrorN())
	}

	// maybe这里还需要进行一次参数校验，一个操作拆成两步，页面上是没问题，但是接口上会有问题，TODO

	d := vo.Data
	standard := vo.Standard

	// 字典类型转换
	dictMapArr := base.GetStructTagsForFieldIndexKey(BomImportVO{}, "dict")
	fields := base.GetStructFields(BomImportVO{})

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
	vos := base.CopyPropertiesList(reflect.TypeOf(model.MoldBom{}), d)
	datas := vos.([]interface{})
	userId := c.GetHeader(constant.USERID)
	var batchSave []model.MoldBom

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)
	// 设置entity对象中的其它字段值
	for _, v := range datas {
		vv := v.(model.MoldBom)
		vv.StandardComponent = standard
		vv.CreatedBy = userId
		vv.UpdatedBy = userId
		batchSave = append(batchSave, vv)
	}

	tx.Table("mold_bom").CreateInBatches(batchSave, len(batchSave))

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 模具BOM
// @Summary 导出
// @Accept json
// @Produce json
// @Param ids query array false "指定ID导出模具BOM"
// @Param standard query string true "是否标准件,Y:标准件 N:非标准件"
// @Param moldCode query string true "模具编码"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold/bom/excel [get]
func ExportBom(c *gin.Context) {
	idsStr := c.QueryArray("ids")
	moldCode := c.QueryArray("moldCode")
	standard := c.QueryArray("standard")

	var res []model.MoldBom
	tx := dao.GetConn().Table("mold_bom").Where("is_deleted = 'N' and mold_code = ? and standard_component = ?", moldCode, standard)

	l := len(idsStr)
	if l == 0 {
		tx.Find(&res)
	} else {
		var ids []int
		for i := 0; i < l; i++ {
			id, err := strconv.Atoi(idsStr[i])
			if err != nil {
				panic(base.ParamsErrorN())
			}

			ids = append(ids, id)
		}

		tx.Where("id in ?", ids).Find(&res)
	}

	vos := base.CopyPropertiesList(reflect.TypeOf(BomExportVO{}), res)

	var vv []BomExportVO
	moldExportVOS := vos.([]interface{})
	for i := 0; i < len(moldExportVOS); i++ {
		vo := moldExportVOS[i].(BomExportVO)
		vv = append(vv, vo)
	}

	fileName := excel.GenerateExcelFile(excel.MOLD_BOM)
	info := common.ExportExcel(fileName, "", vv)

	c.JSON(http.StatusOK, base.Success(info))
}
