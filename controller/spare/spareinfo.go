/**
 * 备件基本信息
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
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// @Tags 备件
// @Summary 添加备件
// @Accept json
// @Produce json
// @Param Body body CreateSpareVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /spare [post]
func CreateSpare(c *gin.Context) {
	var vo CreateSpareVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	// 备件编码唯一
	var count int64
	dao.GetConn().Table("spare_info").Where("code = ? and is_deleted = 'N'", vo.Code).Count(&count)
	if count > 0 {
		panic(base.ResponseEnum[base.SPARE_CODE_EXIST])
	}

	var spareInfo model.SpareInfo
	base.CopyProperties(&spareInfo, vo)

	userId := c.GetHeader(constant.USERID)
	spareInfo.CreatedBy = userId
	spareInfo.UpdatedBy = userId

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	if err := tx.Table("spare_info").Create(&spareInfo).Error; err != nil {
		panic(err)
	}

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(spareInfo.ID))
}

// @Tags 备件
// @Summary 更新备件
// @Accept json
// @Produce json
// @Param Body body UpdateSpareVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /spare [put]
func UpdateSpare(c *gin.Context) {
	var vo UpdateSpareVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	// 存在
	var one model.SpareInfo
	if err := dao.GetConn().Table("spare_info").Where("id = ?", vo.ID).First(&one).Error; err != nil {
		panic(base.ParamsErrorN())
	}

	// 备件编码唯一
	var count int64
	dao.GetConn().Table("spare_info").Where("code = ? and is_deleted = 'N' and id != ?", vo.Code, vo.ID).Count(&count)
	if count > 0 {
		panic(base.ResponseEnum[base.SPARE_CODE_EXIST])
	}

	// 更新
	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	var spareInfo model.SpareInfo
	base.CopyProperties(&spareInfo, vo)
	spareInfo.UpdatedBy = c.GetHeader(constant.USERID)

	if err := tx.Table("spare_info").Where("id = ?", spareInfo.ID).Updates(&spareInfo).Error; err != nil {
		panic(err)
	}

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 备件
// @Summary 用户备件分页
// @Accept json
// @Produce json
// @Param query query PageSpareVO true "PageSpareVO"
// @Param AuthToken header string false "Token"
// @Success 200 {object} spare.PageSpareOutVO
// @Router /spare/page [get]
func PageSpare(c *gin.Context) {
	var vo PageSpareVO
	if err := c.ShouldBindQuery(&vo); err != nil {
		panic(base.ParamsErrorN())
	}

	var spareInfo model.SpareInfo
	base.CopyProperties(&spareInfo, vo)

	tx := dao.GetConn().Table("spare_info")
	if vo.CodeOrName == "" {
		tx = dao.BuildWhereCondition(tx, spareInfo)
	} else {
		tx.Where("(code like concat('%',?,'%') or name like concat('%',?,'%'))", vo.CodeOrName, vo.CodeOrName)
	}

	tx = tx.Where("is_deleted = 'N'").Order("gmt_created desc")

	var result []model.SpareInfo
	page := base.Page(tx, &result, vo.GetCurrentPage(), vo.GetSize())

	// 转换为outvo
	vos := page.List.(*[]model.SpareInfo)
	pageVos, _ := base.CopyPropertiesList(reflect.TypeOf(PageSpareOutVO{}), vos).([]interface{})

	if len(pageVos) == 0 {
		page.List = []interface{}{}
	} else {
		page.List = pageVos
	}

	c.JSON(http.StatusOK, base.Success(page))
}

// @Tags 备件
// @Summary 用户备件查看
// @Accept json
// @Produce json
// @Param id query string true "id"
// @Param AuthToken header string false "Token"
// @Success 200 {object} spare.PageSpareOutVO
// @Router /spare/one [get]
func OneSpare(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		panic(base.ParamsErrorN())
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		panic(base.ParamsErrorN())
	}

	var spareinfo []model.SpareInfo
	dao.GetConn().Table("spare_info").Where("id = ? and is_deleted = 'N'", id).Find(&spareinfo)

	if len(spareinfo) > 0 {
		m := spareinfo[0]
		var vo PageSpareOutVO
		base.CopyProperties(&vo, m)
		c.JSON(http.StatusOK, base.Success(vo))
	} else {
		c.JSON(http.StatusOK, base.SuccessN())
	}
}

// @Tags 备件
// @Summary 删除备件
// @Accept json
// @Produce json
// @Param Body body common.DeleteVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /spare [delete]
func DeleteSpare(c *gin.Context) {
	var vo common.DeleteVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	ids := vo.IDS
	dao.GetConn().Table("spare_info").Where("id in ?", ids).Updates(&model.SpareInfo{
		IsDeleted: "Y",
		UpdatedBy: c.GetHeader(constant.USERID),
	})

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 备件
// @Summary 解析Excel文件内容
// @Accept json
// @Produce json
// @Param file formData file true "文件"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /spare/parse [post]
func ParseSpareExcel(c *gin.Context) {
	header, err := c.FormFile("file")
	if err != nil {
		panic(base.ParamsErrorN())
	}

	// 生成随机excel文件名
	dst := excel.GenerateExcelFile("spare_parse")

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
	vos := common.ParseExcel(dst, SpareImportVO{}, &sb, nil)

	// 对每条数据进行校验
	success := 0
	fault := 0
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
			// 类型强转
			v := vos[i].(SpareImportVO)
			code := v.Code
			var count int64
			dao.GetConn().Table("spare_info").Where("code = ? and is_deleted = 'N'", code).Count(&count)
			if count > 0 {
				m := fmt.Sprintf("第%d行:【备件编号】字段已存在.", i+1)
				sb.WriteString(m)
				fault++
				success--
			}
		}
	}

	// 封装excel头部的所有字段信息
	var head []string
	tags := base.GetStructTagsForFieldIndexKey(SpareImportVO{}, "excel")
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

// @Tags 备件
// @Summary 导入解析后的Excel内容
// @Accept json
// @Produce json
// @Param body body SpareImportParsedVOS true "SpareImportParsedVOS"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /spare/excel [post]
func ImportSpare(c *gin.Context) {
	var vo SpareImportParsedVOS
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	d := vo.Data
	// 字典类型转换
	dictMapArr := base.GetStructTagsForFieldIndexKey(SpareImportVO{}, "dict")
	fields := base.GetStructFields(SpareImportVO{})

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
	vos := base.CopyPropertiesList(reflect.TypeOf(model.SpareInfo{}), d)
	datas := vos.([]interface{})
	userId := c.GetHeader(constant.USERID)
	var batchSave []model.SpareInfo

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)
	// 设置entity对象中的其它字段值
	for _, v := range datas {
		vv := v.(model.SpareInfo)
		vv.CreatedBy = userId
		vv.UpdatedBy = userId
		batchSave = append(batchSave, vv)
	}

	tx.Table("spare_info").CreateInBatches(batchSave, len(batchSave))

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 备件
// @Summary 导出
// @Accept json
// @Produce json
// @Param ids query array false "指定ID导出备件"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /spare/excel [get]
func ExportSpare(c *gin.Context) {
	idsStr := c.QueryArray("ids")

	var res []model.SpareInfo
	tx := dao.GetConn().Table("spare_info")

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
		// dao.BuildWhereIdIn(tx, ids).Find(&res)
	}

	vos := base.CopyPropertiesList(reflect.TypeOf(SpareExportVO{}), res)

	var vv []SpareExportVO
	spareExportVOS := vos.([]interface{})
	for i := 0; i < len(spareExportVOS); i++ {
		vo := spareExportVOS[i].(SpareExportVO)
		vv = append(vv, vo)
	}

	fileName := excel.GenerateExcelFile(excel.SPARE)
	info := common.ExportExcel(fileName, "", vv)

	c.JSON(http.StatusOK, base.Success(info))
}
