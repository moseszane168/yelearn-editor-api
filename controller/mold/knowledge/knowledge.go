/**
 * 模具知识库
 */

package mold

import (
	"crf-mold/base"
	"crf-mold/dao"
	"crf-mold/model"
	"net/http"
	"strconv"

	"crf-mold/common/constant"
	esutil "crf-mold/common/elasticsearch"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// @Tags 知识库
// @Summary 新增知识库
// @Accept json
// @Produce json
// @Param Body body CreateMoldKnowledgeVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /knowledge [post]
func CreateMoldKnowledge(c *gin.Context) {

	var vo CreateMoldKnowledgeVO
	// 传参有误
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	// 标题不重复
	var count int64
	dao.GetConn().Table("mold_knowledge").Where("name = ? and is_deleted = 'N'", vo.Name).Count(&count)
	if count > 0 {
		panic(base.ResponseEnum[base.KNOWLEDGE_EXIST])
	}

	var moldKnowledge model.MoldKnowledge
	base.CopyProperties(&moldKnowledge, vo)

	userId := c.GetHeader(constant.USERID)
	moldKnowledge.CreatedBy = userId
	moldKnowledge.UpdatedBy = userId

	// 开启事务
	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	// 新增
	tx.Table("mold_knowledge").Create(&moldKnowledge)

	// 更新后同步到es
	doc := map[string]interface{}{
		"name":       moldKnowledge.Name,
		"group":      moldKnowledge.Group,
		"content":    esutil.TrimHtml(moldKnowledge.Content),
		"createdBy":  moldKnowledge.CreatedBy,
		"gmtCreated": base.Now(),
	}
	err := esutil.Save(strconv.FormatInt(moldKnowledge.ID, 10), &doc)
	if err != nil {
		panic(base.ResponseEnum[base.ELASTICSEARCH_SAVE_FAULRE])
	}

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(moldKnowledge.ID))
}

// @Tags 知识库
// @Summary 删除知识库
// @Accept json
// @Produce json
// @Param id query int true "id"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /knowledge [delete]
func DeleteMoldKnowledge(c *gin.Context) {
	idstr := c.Query("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		panic(base.ParamsErrorN())
	}

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	tx.Table("mold_knowledge").Where("id = ?", id).Update("is_deleted", "Y")
	if err := esutil.Delete(idstr); err != nil {
		panic(base.ResponseEnum[base.ELASTICSEARCH_SAVE_FAULRE])
	}

	tx.Commit()
	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 知识库
// @Summary 更新知识库
// @Accept json
// @Produce json
// @Param Body body UpdateMoldKnowledgeVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /knowledge [put]
func UpdateMoldKnowledge(c *gin.Context) {
	var vo UpdateMoldKnowledgeVO
	// 传参有误
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	// 资源不存在
	var one model.MoldKnowledge
	if err := dao.GetConn().Table("mold_knowledge").Where("id = ? and is_deleted = 'N'", vo.ID).First(&one).Error; err != nil {
		panic(base.ParamsErrorN())
	}

	// 标题不重复
	var count int64
	dao.GetConn().Table("mold_knowledge").Where("name = ? and is_deleted = 'N' and id != ?", vo.Name, vo.ID).Count(&count)
	if count > 0 {
		panic(base.ResponseEnum[base.KNOWLEDGE_EXIST])
	}

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	var moldKnowledge model.MoldKnowledge
	base.CopyProperties(&moldKnowledge, vo)
	tx.Table("mold_knowledge").Where("id = ?", vo.ID).Updates(&model.MoldKnowledge{
		Name:    vo.Name,
		Content: vo.Content,
		Files:   vo.Files,
		Group:   vo.Group,
		Style:   vo.Style,
	})

	// 更新后同步到es
	doc := map[string]interface{}{
		"name":       moldKnowledge.Name,
		"group":      moldKnowledge.Group,
		"content":    esutil.TrimHtml(moldKnowledge.Content),
		"createdBy":  moldKnowledge.CreatedBy,
		"gmtCreated": base.Now(),
	}
	err := esutil.Save(strconv.FormatInt(moldKnowledge.ID, 10), &doc)
	if err != nil {
		panic(base.ResponseEnum[base.ELASTICSEARCH_SAVE_FAULRE])
	}

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 知识库
// @Summary 知识库列表
// @Accept json
// @Produce json
// @Param name query string false "知识库标题"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /knowledge [get]
func ListMoldKnowledge(c *gin.Context) {
	name := c.Query("name")

	var results []model.MoldKnowledge
	tx := dao.GetConn().Table("mold_knowledge").Where("is_deleted = 'N'")
	if name != "" {
		tx.Where("name like concat('%',?,'%')", name)
	}

	tx.Order("`gmt_created` desc").Find(&results)

	if len(results) > 0 {
		c.JSON(http.StatusOK, base.Success(results))
	} else {
		c.JSON(http.StatusOK, base.Success([]model.MoldKnowledge{}))
	}
}

// @Tags 知识库
// @Summary 知识库分页
// @Accept json
// @Produce json
// @Param name query string false "知识库标题"
// @Param currentPage query int false "当前页"
// @Param size query int false "每页数量"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /knowledge/page [get]
func PageMoldKnowledge(c *gin.Context) {
	name := c.Query("name")

	current, size := base.GetPageParams(c)
	var results []model.MoldKnowledge
	tx := dao.GetConn().Table("mold_knowledge").Where("is_deleted = 'N'").Order("`gmt_created` desc")
	if name != "" {
		tx.Where("name like concat('%',?,'%')", name)
	}

	page := base.Page(tx, &results, current, size)

	if len(results) == 0 {
		page.List = []interface{}{}
	}

	c.JSON(http.StatusOK, base.Success(page))
}

// @Tags 知识库
// @Summary 知识库查看
// @Accept json
// @Produce json
// @Param id query string true "id"
// @Param AuthToken header string false "Token"
// @Router /knowledge/one [get]
func OneKnowledge(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		panic(base.ParamsErrorN())
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		panic(base.ParamsErrorN())
	}

	var knowledge []model.MoldKnowledge
	dao.GetConn().Table("mold_knowledge").Where("id = ? and is_deleted = 'N'", id).Find(&knowledge)

	if len(knowledge) > 0 {
		c.JSON(http.StatusOK, base.Success(knowledge[0]))
	} else {
		c.JSON(http.StatusOK, base.SuccessN())
	}
}

// @Tags 知识库
// @Summary 知识库搜索
// @Accept json
// @Produce json
// @Param q query string true "q"
// @Param currentPage query int false "当前页"
// @Param size query int false "每页数量"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /knowledge/search [get]
func SearchKnowledge(c *gin.Context) {
	q := c.Query("q")
	current, size := base.GetPageParams(c)

	query := map[string]interface{}{
		"_source": []string{"group", "gmtCreated", "createdBy", "name"},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"name": q,
						},
					},
					{
						"match": map[string]interface{}{
							"content": q,
						},
					},
				},
			},
		},
		"highlight": map[string]interface{}{
			"pre_tags":  []string{"<em>"},
			"post_tags": []string{"</em>"},
			"fields": map[string]interface{}{
				"name":    map[string]interface{}{},
				"content": map[string]interface{}{},
			},
		},
		"from": (current - 1) * size,
		"size": size,
	}

	r, err := esutil.Search(&query)
	if err != nil {
		panic(err)
	}

	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
	total := r["hits"].(map[string]interface{})["total"].(map[string]interface{})
	totalRow := int64(total["value"].(float64))

	result := make([]KnowledgeSearchVO, len(hits))

	for i := 0; i < len(hits); i++ {
		item := hits[i]
		highlight := item.(map[string]interface{})["highlight"].((map[string]interface{}))
		source := item.(map[string]interface{})["_source"].((map[string]interface{}))

		// id
		idstr := item.(map[string]interface{})["_id"].(string)
		id, _ := strconv.ParseInt(idstr, 10, 64)

		// name,标题没有则取原始数据标题
		name := ""

		nameinf, ok := highlight["name"]
		if ok {
			namearr, ok := nameinf.([]interface{})
			if ok {
				name = namearr[0].(string)
			}
		} else {
			name = source["name"].(string)
		}

		content := ""
		contentinf, ok := highlight["content"]
		if ok {
			contentarr, ok := contentinf.([]interface{})
			if ok {
				content = contentarr[0].(string)
			}
		} else {
			content = source["content"].(string)
		}

		// createby
		createBy := source["createdBy"].(string)

		// gmtcreated
		gmtCreated := source["gmtCreated"].(string)

		// group
		group := source["group"].(string)

		result[i].ID = id
		result[i].Content = content
		result[i].CreatedBy = createBy
		result[i].GmtUpdated = gmtCreated
		result[i].Group = group
		result[i].Name = name
	}

	page := base.NewBasePage(current, size, totalRow, result)
	if len(result) == 0 {
		page.List = []interface{}{}
	} else {
		page.List = result
	}

	c.JSON(http.StatusOK, base.Success(page))
}
