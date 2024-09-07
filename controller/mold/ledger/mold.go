/**
 * 模具台账
 */

package mold

import (
	"crf-mold/base"
	"crf-mold/common"
	"crf-mold/common/constant"
	"crf-mold/common/validate"
	"crf-mold/controller/dict"
	"crf-mold/controller/excel"
	"crf-mold/controller/user"
	"crf-mold/dao"
	"crf-mold/model"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// @Tags 台账
// @Summary 添加台账
// @Accept json
// @Produce json
// @Param Body body CreateMoldVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold [post]
func CreateMold(c *gin.Context) {
	var vo CreateMoldVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	// 模具编码唯一
	var count int64
	dao.GetConn().Table("mold_info").Where("code = ? and is_deleted = 'N'", vo.Code).Count(&count)
	if count > 0 {
		panic(base.ResponseEnum[base.MOLD_CODE_EXIST])
	}

	var moldInfo model.MoldInfo
	base.CopyProperties(&moldInfo, vo)

	userId := c.GetHeader(constant.USERID)
	moldInfo.CreatedBy = userId
	moldInfo.UpdatedBy = userId

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	if err := tx.Table("mold_info").Create(&moldInfo).Error; err != nil {
		panic(err)
	}

	if len(vo.Customs) > 0 {
		if err := SaveCustomInfos(tx, vo.Customs, moldInfo.Code, moldInfo.Code, c.GetHeader(constant.USERID)); err != nil {
			panic(err)
		}
	}

	if len(vo.PartCodes) > 0 {
		if err := SavePartCodes(tx, vo.PartCodes, moldInfo.Code, moldInfo.Code, c.GetHeader(constant.USERID)); err != nil {
			panic(err)
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(moldInfo.ID))
}

// @Tags 台账
// @Summary 更新台账
// @Accept json
// @Produce json
// @Param Body body UpdateMoldVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold [put]
func UpdateMold(c *gin.Context) {
	var vo UpdateMoldVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	// 存在
	var one model.MoldInfo
	if err := dao.GetConn().Table("mold_info").Where("id = ?", vo.ID).First(&one).Error; err != nil {
		panic(base.ParamsErrorN())
	}

	// 模具编码唯一
	var count int64
	dao.GetConn().Table("mold_info").Where("code = ? and is_deleted = 'N' and id != ?", vo.Code, vo.ID).Count(&count)
	if count > 0 {
		panic(base.ResponseEnum[base.MOLD_CODE_EXIST])
	}

	// 更新
	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	// 如果更新模具编码，将会修改所有通过模具编码关联的其它表信息
	if vo.Code != one.Code {
		tx.Table("mold_bom").Where("mold_code = ?", one.Code).Update("mold_code", vo.Code)
		tx.Table("mold_doc").Where("mold_code = ?", one.Code).Update("mold_code", vo.Code)
		tx.Table("mold_params").Where("mold_code = ?", one.Code).Update("mold_code", vo.Code)
		tx.Table("mold_params_custom").Where("mold_code = ?", one.Code).Update("mold_code", vo.Code)
		tx.Table("mold_remodel").Where("mold_code = ?", one.Code).Update("mold_code", vo.Code)
		tx.Table("mold_repair").Where("mold_code = ?", one.Code).Update("mold_code", vo.Code)
	}

	// 更新模具
	var moldInfo model.MoldInfo
	base.CopyProperties(&moldInfo, vo)
	moldInfo.UpdatedBy = c.GetHeader(constant.USERID)

	if err := tx.Table("mold_info").Where("id = ?", moldInfo.ID).Updates(&moldInfo).Error; err != nil {
		panic(err)
	}
	// 针对RFID空字符串的更新
	if err := tx.Table("mold_info").Where("id = ?", moldInfo.ID).Updates(map[string]interface{}{"rfid": vo.Rfid}).Error; err != nil {
		panic(err)
	}

	if len(vo.Customs) > 0 {
		if err := SaveCustomInfos(tx, vo.Customs, one.Code, moldInfo.Code, c.GetHeader(constant.USERID)); err != nil {
			panic(err)
		}
	}

	if len(vo.PartCodes) > 0 {
		if err := SavePartCodes(tx, vo.PartCodes, one.Code, moldInfo.Code, c.GetHeader(constant.USERID)); err != nil {
			panic(err)
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 台账
// @Summary 用户台账分页
// @Accept json
// @Produce json
// @Param query query PageMoldVO true "PageMoldVO"
// @Param AuthToken header string false "Token"
// @Success 200 {object} mold.PageMoldOutV1
// @Router /mold/page [get]
func PageMold(c *gin.Context) {
	var vo PageMoldVO
	if err := c.ShouldBindQuery(&vo); err != nil {
		panic(base.ParamsError(err.Error()))
	}
	var page *base.BasePage
	//var result []PageMoldOutVO
	var result []PageMoldOutV1
	sql := `
		select 
			m.id,
			m.code,
			m.name,
			m.client_name,
			m.weight,
			m.project_name,
			m.size_long,
			m.size_width,
			m.size_heigh,
			m.platform,
			m.process,
			m.property_number,
			m.type,
			m.status,
			m.make_date,
			m.first_use_date,
			m.year_limit,
			m.rfid,
			m.line_level,
			m.provide,
			m.flush_count,
			m.category,
			GROUP_CONCAT(rel.part_code SEPARATOR '/') as part_code_str,
			CASE WHEN p.id IS NULL THEN false ELSE true END AS have_plan
		from 
			mold_info m
        left join mold_part_rel rel on 
			m.code = rel.mold_code and rel.is_deleted = 'N'
		left join mold_maintenance_plan_rel p on m.id = p.mold_id and p.is_deleted = 'N'
		where
			m.is_deleted = 'N'
	`
	if vo.ExcludeType != "" {
		combineString := "'" + strings.Join(strings.Split(vo.ExcludeType, ","), "','") + "'"
		sql += "and m.type not in (" + combineString + ")"
	}
	if vo.CodeOrName == "" {
		// 额外处理零件号字段
		if vo.PartCodes != "" {
			sql += fmt.Sprintf(" and rel.part_code like concat('%%', '%s','%%')", vo.PartCodes)
		}
		if vo.Code != "" {
			sql += fmt.Sprintf(" and m.code like concat('%%', '%s','%%')", vo.Code)
		}
		if vo.Name != "" {
			sql += fmt.Sprintf(" and m.name like concat('%%', '%s','%%')", vo.Name)
		}
		// 额外处理时间字段
		if vo.FirstUseDateBegin != nil {
			sql += fmt.Sprintf(" and m.first_use_date >= '%s'", vo.FirstUseDateBegin)
		}
		if vo.FirstUseDateEnd != nil {
			sql += fmt.Sprintf(" and m.first_use_date <= '%s'", vo.FirstUseDateEnd)
		}
		if vo.MakeDateBegin != nil {
			sql += fmt.Sprintf(" and m.make_date >= '%s'", vo.MakeDateBegin)
		}
		if vo.MakeDateEnd != nil {
			sql += fmt.Sprintf(" and m.make_date <= '%s'", vo.MakeDateEnd)
		}
		if vo.Process != "" {
			sql += fmt.Sprintf(" and m.process = '%s'", vo.Process)
		}
		if vo.Type != "" {
			sql += fmt.Sprintf(" and m.type = '%s'", vo.Type)
		}
		if vo.Status != "" {
			sql += fmt.Sprintf(" and m.status = '%s'", vo.Status)
		}
		if vo.LineLevel != "" {
			sql += fmt.Sprintf(" and m.line_level = '%s'", vo.LineLevel)
		}
		if vo.Platform != "" {
			sql += fmt.Sprintf(" and m.platform = '%s'", vo.Platform)
		}
		if vo.Rfid != "" {
			sql += fmt.Sprintf(" and m.rfid = '%s'", vo.Rfid)
		}
		if vo.ProjectName != "" {
			sql += fmt.Sprintf(" and m.project_name like concat('%%', '%s','%%')", vo.ProjectName)
		}
	} else {
		sql += fmt.Sprintf(" and (m.code like concat('%%', '%s','%%') or m.name like concat('%%', '%s','%%') or rel.part_code like concat('%%', '%s','%%'))", vo.CodeOrName, vo.CodeOrName, vo.CodeOrName)
	}

	sql += fmt.Sprintf("group by m.id order by m.gmt_created desc")

	var tx = dao.GetConn()
	var countVo base.CountVO
	current, size := vo.GetCurrentPage(), vo.GetSize()
	countSql := "select count(1) as count from (" + sql + ") t"
	tx.Raw(countSql).Scan(&countVo)

	pageSql := sql + fmt.Sprintf(" limit %d,%d", int((current-1)*size), size)
	tx.Raw(pageSql).Scan(&result)

	for i, row := range result {
		if row.PartCodeStr != "" {
			result[i].PartCodes = strings.Split(row.PartCodeStr, "/")
		} else {
			result[i].PartCodes = []string{}
		}
	}
	page = base.NewBasePage(current, size, countVo.Count, result)
	if len(result) == 0 {
		page.List = []interface{}{}
	} else {
		page.List = result
	}

	c.JSON(http.StatusOK, base.Success(page))
}

// @Tags 台账
// @Summary 用户台账查看
// @Accept json
// @Produce json
// @Param id query string true "id"
// @Param AuthToken header string false "Token"
// @Success 200 {object} mold.PageMoldOutVO
// @Router /mold/one [get]
func OneMold(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		panic(base.ParamsErrorN())
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		panic(base.ParamsErrorN())
	}

	var moldinfo []model.MoldInfo
	dao.GetConn().Table("mold_info").Where("id = ? and is_deleted = 'N'", id).Find(&moldinfo)

	if len(moldinfo) > 0 {
		m := moldinfo[0]
		var vo PageMoldOutVO
		base.CopyProperties(&vo, m)

		// 零件号
		partCodes := GetMoldPartCodes(m.Code)
		if partCodes == nil {
			vo.PartCodes = []string{}
		} else {
			vo.PartCodes = partCodes
		}

		// 自定义字段
		var moldCustoms []model.MoldCustomInfo
		dao.GetConn().Table("mold_custom_info").Where("mold_code = ? and is_deleted = 'N'", m.Code).Find(&moldCustoms)
		if len(moldCustoms) == 0 {
			vo.Customs = []MoldCustomVO{}
		} else {
			moldCustomVos := base.CopyPropertiesList(reflect.TypeOf(MoldCustomVO{}), moldCustoms)
			vo.Customs = moldCustomVos
		}

		c.JSON(http.StatusOK, base.Success(vo))
	} else {
		c.JSON(http.StatusOK, base.SuccessN())
	}
}

// @Tags 台账
// @Summary 删除台账
// @Accept json
// @Produce json
// @Param Body body common.DeleteVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold [delete]
func DeleteMold(c *gin.Context) {
	var vo common.DeleteVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	ids := vo.IDS
	// 查出所有模具的code
	var codes []string
	var result []string

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	tx.Table("mold_info").Where("id in ?", ids).Select("DISTINCT(code)").Find(&codes)
	tx.Table("mold_stereoscopic_warehouse_location").Select("mold_code").Where(
		"mold_code in ?", codes).Scan(&result)
	if len(result) > 0 {
		response := base.Response{Code: base.MOLD_EXIST_IN_ASRS, Message: fmt.Sprintf("模具 %s 存在立库中, 不允许删除", result[0]), Result: result[0]}
		panic(&response)
	}

	// 删除关联表
	tx.Table("mold_bom").Where("mold_code in ?", codes).Update("is_deleted", "Y")
	tx.Table("mold_custom_info").Where("mold_code in ?", codes).Update("is_deleted", "Y")
	tx.Table("mold_doc").Where("mold_code in ?", codes).Update("is_deleted", "Y")
	tx.Table("mold_params").Where("mold_code in ?", codes).Update("is_deleted", "Y")
	tx.Table("mold_params_custom").Where("mold_code in ?", codes).Update("is_deleted", "Y")
	tx.Table("mold_part_rel").Where("mold_code in ?", codes).Update("is_deleted", "Y")
	tx.Table("mold_remodel").Where("mold_code in ?", codes).Update("is_deleted", "Y")
	tx.Table("mold_repair").Where("mold_code in ?", codes).Update("is_deleted", "Y")

	tx.Table("mold_info").Where("id in ?", ids).Updates(&model.MoldInfo{
		IsDeleted: "Y",
		UpdatedBy: c.GetHeader(constant.USERID),
	})

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 台账
// @Summary 解析Excel文件内容
// @Accept json
// @Produce json
// @Param file formData file true "文件"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold/parse [post]
func ParseMoldExcel(c *gin.Context) {
	header, err := c.FormFile("file")
	if err != nil {
		panic(base.ParamsErrorN())
	}

	// 生成随机excel文件名
	dst := excel.GenerateExcelFile("mold_parse")

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

	var sb strings.Builder
	// 解析上传的excel文件
	vos := common.ParseExcel(dst, MoldImportVO{}, &sb, nil)

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
			v := vos[i].(MoldImportVO)
			code := v.Code
			var count int64
			dao.GetConn().Table("mold_info").Where("code = ? and is_deleted = 'N'", code).Count(&count)
			if count > 0 {
				m := fmt.Sprintf("第%d行:【模具编号】字段已存在.", i+1)
				sb.WriteString(m)
				fault++
				success--
			}
		}
	}

	// 封装excel头部的所有字段信息
	var head []string
	tags := base.GetStructTagsForFieldIndexKey(MoldImportVO{}, "excel")
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

// @Tags 台账
// @Summary 导入解析后的Excel内容
// @Accept json
// @Produce json
// @Param body body MoldImportParsedVOS true "MoldImportParsedVOS"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold/excel [post]
func ImportMold(c *gin.Context) {
	var vo MoldImportParsedVOS
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	// maybe这里还需要进行一次参数校验，一个操作拆成两步，页面上是没问题，但是接口上会有问题，TODO

	d := vo.Data

	// 字典类型转换
	dictMapArr := base.GetStructTagsForFieldIndexKey(MoldImportVO{}, "dict")
	fields := base.GetStructFields(MoldImportVO{})

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
	vos := base.CopyPropertiesList(reflect.TypeOf(model.MoldInfo{}), d)
	datas := vos.([]interface{})
	userId := c.GetHeader(constant.USERID)
	var batchSave []model.MoldInfo

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)
	// 设置entity对象中的其它字段值
	for i, v := range datas {
		vv := v.(model.MoldInfo)
		vv.CreatedBy = userId
		vv.UpdatedBy = userId

		// 使用年限计算，设置
		if vv.FirstUseDate != nil {
			beginDate := *vv.FirstUseDate

			yearLimit := float64(time.Now().Unix()-beginDate.Time().Unix()) / 60 / 60.0 / 24.0 / 365.0
			vv.YearLimit = yearLimit
		}

		// 保存零件号
		vo := d[i]
		partCodesStr := vo.PartCodes
		if partCodesStr != "" {
			partCodes := strings.Split(partCodesStr, "/")
			if err := SavePartCodes(tx, partCodes, vv.Code, vv.Code, userId); err != nil {
				panic(err)
			}
		}

		batchSave = append(batchSave, vv)
	}

	tx.Table("mold_info").CreateInBatches(batchSave, len(batchSave))

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 台账
// @Summary 导出
// @Accept json
// @Produce json
// @Param ids query array false "指定ID导出台账"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold/excel [get]
func ExportMold(c *gin.Context) {
	idsStr := c.QueryArray("ids")

	var res []model.MoldInfo
	tx := dao.GetConn().Table("mold_info").Where("is_deleted = 'N'")

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

	vos := base.CopyPropertiesList(reflect.TypeOf(MoldExportVO{}), res)

	var vv []MoldExportVO
	moldExportVOS := vos.([]interface{})
	for i := 0; i < len(moldExportVOS); i++ {
		vo := moldExportVOS[i].(MoldExportVO)
		// 查询零件号
		var sb strings.Builder
		var partRels []model.MoldPartRel
		dao.GetConn().Table("mold_part_rel").Where("mold_code = ? and is_deleted = 'N'", vo.Code).Find(&partRels)
		for _, v := range partRels {
			sb.WriteString(v.PartCode)
			sb.WriteString(",")
		}
		partCodesStr := sb.String()
		l := len(partCodesStr)
		if l > 0 {
			vo.PartCodes = partCodesStr[:l-1]
		}
		vv = append(vv, vo)
	}

	fileName := excel.GenerateExcelFile(excel.MOLD)
	info := common.ExportExcel(fileName, "", vv)

	c.JSON(http.StatusOK, base.Success(info))
}

// @Tags 台账
// @Summary 置零模具使用冲次
// @Accept json
// @Produce json
// @Param Body body MoldFlushCountVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold/flush [put]
func ZeroMoldFlushCount(c *gin.Context) {
	var vo MoldFlushCountVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	// 管理员操作
	userId := c.GetHeader(constant.USERID)
	if userId != "admin" {
		panic(base.ResponseEnum[base.MOLD_FLUSH_COUNT_ZERO_ONLY_ADMIN])
	}

	// 用户不存在
	var u model.UserInfo
	if err := dao.GetConn().Table("user_info").Where("login_name = ? and is_deleted = 'N'", userId).First(&u).Error; err != nil {
		panic(base.ParamsErrorN())
	}

	// 解密旧密码
	old, err := user.DecrtptPassword(vo.PassWord)
	if err != nil {
		panic(base.ParamsErrorN())
	}

	// 旧密码是否正确
	md5Pwd := base.MD5(string(old))
	if u.Password != md5Pwd {
		panic(base.ResponseEnum[base.USER_OR_PASSWORD_ERROR])
	}

	dao.GetConn().Table("mold_info").Where("id = ? and is_deleted = 'N'", vo.ID).Updates(map[string]interface{}{
		"flush_count":      0,
		"calc_flush_count": 0,
		"updated_by":       userId,
	})

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 台账
// @Summary 查询模具所有零件号
// @Accept json
// @Produce json
// @Param AuthToken header string false "Token"
// @Success 200 {string} json string
// @Router /mold/part/code/group [get]
func MoldPartCodeGroup(c *gin.Context) {
	var result []string
	dao.GetConn().Table("mold_part_rel").Select("part_code").Where("is_deleted ='N'").Group("part_code").Find(&result)
	c.JSON(http.StatusOK, base.Success(result))
}

// @Tags 台账
// @Summary 查询模具编码根据零件号
// @Accept json
// @Produce json
// @Param Body body PartCodeQueryMoldVO true "PartCodeQueryMoldVO"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json string
// @Router /mold/by/part/code [post]
func MoldQueryByPartCode(c *gin.Context) {
	var vo PartCodeQueryMoldVO

	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	var result []string
	dao.GetConn().Table("mold_part_rel").Select("mold_code").Where("is_deleted ='N' and part_code = ?", vo.PartCode).Group("mold_code").Find(&result)
	c.JSON(http.StatusOK, base.Success(result))
}
