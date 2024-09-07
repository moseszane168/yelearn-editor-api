package coursepacks

import (
	"bytes"
	"crf-mold/base"
	"crf-mold/common/constant"
	"crf-mold/dao"
	"crf-mold/model/earthworm"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io"
	"net/http"
	"strconv"
)

// @Tags 课程包
// @Summary 新增课程包
// @Accept json
// @Produce json
// @Param Body body CreateCoursePackVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {object} bool
// @Router /editor/course-pack [post]
func CreateCoursePack(c *gin.Context) {
	var vo CreateCoursePackVO
	// 传参校验
	// 在检验前先保存请求体的内容 以便后续再次校验
	bodyBytes, _ := io.ReadAll(c.Request.Body)
	// 重新设置请求体的内容 否则ShouldBindBodyWith无法读取到
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	// 执行校验
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	// 名称唯一
	/*var count int64
	dao.GetConn().Table("course_packs").
		Where("title = ?", vo.Name).Count(&count)
	if count > 0 {
		panic(base.ResponseEnum[base.MAINTENANCE_PLAN_NAME_EXIST])
	}*/

	userId := c.GetHeader(constant.USERID)

	var coursePack earthworm.CoursePacks
	base.CopyProperties(&coursePack, vo)
	coursePack.IsFree = true
	coursePack.Order = 1
	//coursePack.CreatedAt = CreatedAt
	//coursePack.UpdatedAt = UpdatedAt
	coursePack.CreatorId = userId
	coursePack.ShareLevel = "public"
	coursePack.Cover = vo.Cover
	coursePack.Title = vo.Title
	coursePack.Description = vo.Description
	coursePack.Id = base.GenerateUniqueTextID()

	// 新增
	dao.GetConn().Table("course_packs").Create(&coursePack)

	c.JSON(http.StatusOK, base.Success(coursePack.Id))
}

// @Tags 课程包
// @Summary 删除课程包
// @Accept json
// @Produce json
// @Param id query int true "id"
// @Param AuthToken header string false "Token"
// @Success 200 {object} bool
// @Router /editor/course-pack [delete]
func DeleteCoursePack(c *gin.Context) {
	idstr := c.Query("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		panic(base.ParamsErrorN())
	}

	//userId := c.GetHeader(constant.USERID)

	var coursePack earthworm.CoursePacks
	dao.GetConn().Table("course_packs").
		Where("id = ?", id).
		Delete(&coursePack)

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 课程包
// @Summary 更新课程包
// @Accept json
// @Produce json
// @Param Body body UpdateCoursePackVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {object} bool
// @Router /editor/course-pack [put]
func UpdateCoursePack(c *gin.Context) {
	var vo UpdateCoursePackVO
	// 传参有误
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	//userId := c.GetHeader(constant.USERID)

	// 名称唯一
	/*var c1 int64
	dao.GetConn().Table("mold_maintenance_plan").
		Where("name = ? and id != ? and is_deleted = 'N'", vo.Name, vo.ID).Count(&c1)
	if c1 > 0 {
		panic(base.ResponseEnum[base.MAINTENANCE_NAME_EXIST])
	}*/

	var coursePack earthworm.CoursePacks
	base.CopyProperties(&coursePack, vo)

	dao.GetConn().Table("course_packs").
		Updates(&coursePack)

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 课程包
// @Summary 课程包列表
// @Accept json
// @Produce json
// @Param id query int true "id"
// @Param AuthToken header string false "Token"
// @Success 200 {object} earthworm.CoursePacks
// @Router /editor/course-packs [get]
func ListCoursePack(c *gin.Context) {
	id := c.Query("id")

	var results []earthworm.CoursePacks
	tx := dao.GetConn().Table("course_packs").
		Where("id = ?", id)

	tx.Order("`gmt_created` desc").Find(&results)

	if len(results) > 0 {
		c.JSON(http.StatusOK, base.Success(results))
	} else {
		c.JSON(http.StatusOK, base.Success([]earthworm.CoursePacks{}))
	}
}

// @Tags 课程包
// @Summary 课程包分页
// @Accept json
// @Produce json
// @Param query query PageMaintenancePlanInputVO true "PageMaintenancePlanInputVO"
// @Param AuthToken header string false "Token"
// @Success 200 {object} PageMaintenancePlanOutputVO
// @Router /maintenance/plan/page [get]
func PageCoursePack(c *gin.Context) {
	/*var vo PageMaintenancePlanInputVO
	if err := c.ShouldBindQuery(&vo); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	current, size := base.GetPageParams(c)
	var results []model.MoldMaintenancePlan
	tx := dao.GetConn().Table("mold_maintenance_plan").Where("is_deleted = 'N'").Order("`gmt_created` desc")
	if vo.CodeOrName != "" {
		tx.Where("(code like concat('%',?,'%') or name like concat('%',?,'%'))", vo.CodeOrName, vo.CodeOrName)
	} else {
		m := model.MoldMaintenancePlan{}
		base.CopyProperties(&m, vo)
		dao.BuildWhereCondition(tx, m)
		// 额外处理时间字段
		if vo.GmtCreatedBegin != "" {
			tx.Where("gmt_created >= ?", vo.GmtCreatedBegin)
		}
		if vo.GmtCreatedEnd != "" {
			tx.Where("gmt_created <= ?", vo.GmtCreatedEnd)
		}
	}

	page := base.Page(tx, &results, current, size)

	if len(results) == 0 {
		page.List = []interface{}{}
	} else {
		outVos := make([]PageMaintenancePlanOutputVO, len(results))
		list := page.List.(*[]model.MoldMaintenancePlan)
		for i := 0; i < len(*list); i++ {
			base.CopyProperties(&(outVos[i]), (*list)[i])

			// 设置一下部门
			createdBy := outVos[i].CreatedBy
			if createdBy != "" {
				var ui model.UserInfo
				if err := dao.GetConn().Table("user_info").Where("login_name = ?", outVos[i].CreatedBy).First(&ui).Error; err == nil {
					outVos[i].Department = ui.Department
					outVos[i].CreatedBy = ui.Name
				}
			}
		}

		page.List = outVos
	}*/

	//c.JSON(http.StatusOK, base.Success(page))
}

// @Tags 课程包
// @Summary 课程包查看
// @Accept json
// @Produce json
// @Param id query string true "id"
// @Param AuthToken header string false "Token"
// @Success 200 {object} earthworm.CoursePacks
// @Router /editor/course-pack [get]
func OneCoursePack(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		panic(base.ParamsErrorN())
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		panic(base.ParamsErrorN())
	}

	var coursePack earthworm.CoursePacks
	if err := dao.GetConn().Table("course_packs").
		Where("id = ?", id).
		First(&coursePack).Error; err != nil {
		panic(base.ParamsErrorN())
	}

	c.JSON(http.StatusOK, base.Success(coursePack))
}
