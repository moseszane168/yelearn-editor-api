package statements

import (
	"bytes"
	"crf-mold/base"
	"crf-mold/dao"
	"crf-mold/model/earthworm"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strings"
)

// @Tags 句子
// @Summary 新增句子
// @Accept json
// @Produce json
// @Param Body body CreateStatementVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {object} bool
// @Router /editor/statement [post]
func CreateStatement(c *gin.Context) {
	var vo CreateStatementVO
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
	dao.GetConn().Table("courses").
		Where("title = ?", vo.Name).Count(&count)
	if count > 0 {
		panic(base.ResponseEnum[base.MAINTENANCE_PLAN_NAME_EXIST])
	}*/

	//userId := c.GetHeader(constant.USERID)

	var statement earthworm.Statements
	base.CopyProperties(&statement, vo)
	statement.Order = 1
	//statement.CreatedAt = CreatedAt
	//statement.UpdatedAt = UpdatedAt
	statement.CourseId = vo.CourseId
	//statement.Soundmark = vo.SoundMark
	//statement.Chinese = vo.Chinese
	statement.English = vo.English
	statement.Id = base.GenerateUniqueTextID()

	// 新增
	errCreate := dao.GetConn().Table("statements").
		Create(&statement).Error
	if errCreate != nil {
		panic(base.ParamsError(errCreate.Error()))
	}

	c.JSON(http.StatusOK, base.Success(statement.Id))
}

// @Tags 句子
// @Summary 删除句子
// @Accept json
// @Produce json
// @Param id query int true "id"
// @Param AuthToken header string false "Token"
// @Success 200 {object} bool
// @Router /editor/statement [delete]
func DeleteStatement(c *gin.Context) {
	var vo DeleteStatementVO
	// 传参有误
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	var statement earthworm.Statements
	if len(vo.StatementIds) == 0 {
		errDelete := dao.GetConn().Table("statements").
			Where("course_id = ?", vo.CourseId).
			Delete(&statement).Error
		if errDelete != nil {
			panic(base.ParamsError(errDelete.Error()))
		}
	} else {
		errDelete := dao.GetConn().Table("statements").
			Where("course_id = ?", vo.CourseId).
			Where("id in ?", vo.StatementIds).
			Delete(&statement).Error
		if errDelete != nil {
			panic(base.ParamsError(errDelete.Error()))
		}
	}

	//userId := c.GetHeader(constant.USERID)

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 句子
// @Summary 更新句子
// @Accept json
// @Produce json
// @Param Body body UpdateStatementVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {object} bool
// @Router /editor/statement [put]
func UpdateStatement(c *gin.Context) {
	var vo UpdateStatementVO
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

	var statement earthworm.Statements
	base.CopyProperties(&statement, vo)

	errUpdates := dao.GetConn().Table("statements").
		Where("id = ?", vo.Id).
		Updates(&statement).Error
	if errUpdates != nil {
		panic(base.ParamsError(errUpdates.Error()))
	}

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 句子
// @Summary 句子列表
// @Accept json
// @Produce json
// @Param id query int true "id"
// @Param AuthToken header string false "Token"
// @Success 200 {object} earthworm.Statements
// @Router /editor/statements [get]
func ListStatement(c *gin.Context) {
	var vo ListStatementVO
	// 传参有误
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	var statements []earthworm.Statements
	statementMap := make(map[string][]earthworm.Statements)
	if vo.StatementId == "statement" {
		errFind := dao.GetConn().Table("statements").
			Where("course_id = ?", vo.CourseId).
			Where("pid <> ?", "").
			//Order("`order` desc").
			Find(&statements).Error
		if errFind != nil {
			panic(base.ParamsError(errFind.Error()))
		}
		for _, statement := range statements {
			statementMap[statement.Pid] = append(statementMap[statement.Pid], statement)
		}

		if len(statementMap) > 0 {
			c.JSON(http.StatusOK, base.Success(statementMap))
		} else {
			c.JSON(http.StatusOK, base.Success(nil))
		}
	} else {
		errFind := dao.GetConn().Table("statements").
			Where("course_id = ?", vo.CourseId).
			Where("pid = ?", "").
			//Order("`order` desc").
			Find(&statements).Error
		if errFind != nil {
			panic(base.ParamsError(errFind.Error()))
		}

		if len(statements) > 0 {
			c.JSON(http.StatusOK, base.Success(statements))
		} else {
			c.JSON(http.StatusOK, base.Success([]earthworm.Courses{}))
		}
	}
}

// @Tags 句子
// @Summary 句子分页
// @Accept json
// @Produce json
// @Param query query PageMaintenancePlanInputVO true "PageMaintenancePlanInputVO"
// @Param AuthToken header string false "Token"
// @Success 200 {object} PageMaintenancePlanOutputVO
// @Router /maintenance/plan/page [get]
func PageStatement(c *gin.Context) {
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

// @Tags 句子
// @Summary 句子查看
// @Accept json
// @Produce json
// @Param id query string true "id"
// @Param AuthToken header string false "Token"
// @Success 200 {object} earthworm.Statements
// @Router /editor/course-pack [get]
func OneStatement(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		panic(base.ParamsErrorN())
	}

	var statement earthworm.Statements
	errFirst := dao.GetConn().Table("statements").
		Where("id = ?", id).
		First(&statement).Error
	if errFirst != nil {
		panic(base.ParamsError(errFirst.Error()))
	}

	// todo 课程详情

	// todo 课程包详情

	c.JSON(http.StatusOK, base.Success(statement))
}

// @Tags 句子
// @Summary 拆分句子
// @Accept json
// @Produce json
// @Param id query string true "id"
// @Param AuthToken header string false "Token"
// @Success 200 {object} earthworm.Statements
// @Router /editor/split-statement [post]
func SplitStatement(c *gin.Context) {
	var vo SplitStatementVO
	// 传参有误
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	/*var sentences []earthworm.Statements
	if len(vo.StatementIds) == 0 {
		// todo 查出该课程所有句子
		errFind := dao.GetConn().Table("statements").
			Where("course_id = ?", vo.CourseId).
			Where("pid = ?", "").
			//Order("`order` desc").
			Find(&sentences).Error
		if errFind != nil {
			panic(base.ParamsError(errFind.Error()))
		}
	} else {
		// todo 查出该句子
		errFind := dao.GetConn().Table("statements").
			Where("course_id = ?", vo.CourseId).
			Where("id in ?", vo.StatementIds).
			//Order("`order` desc").
			Find(&sentences).Error
		if errFind != nil {
			panic(base.ParamsError(errFind.Error()))
		}
	}

	if len(sentences) == 0 {
		panic(base.ParamsError("请添加句子"))
	}

	var sentencesSlice []string
	for i, _ := range sentences {
		sentencesSlice = append(sentencesSlice, sentences[i].English+".")
	}
	sentencesStr := strings.Join(sentencesSlice, "")
	//sentencesStr = sentencesStr + "Come here and learn English quickly."

	content := "请注意，我将给你发送一些英文句子，你需要将英文句子拆分成单词和短语句子，按照以下顺序拆分：" +
		"1. 拆分出句子中第一个短句里的每个单词。" +
		"2. 将第一个短句的所有单词组合成一个短句。" +
		"3. 继续拆分下一个短句的所有单词。" +
		"4. 将下一个短句的所有单词组合成一个短句。" +
		"5. 将前两个短句组合成一个更长的短句。" +
		"6. 继续此过程，直到所有短句都被拆分并组合成短句。" +
		"7. 注意，组成短句要合理。例如：so i don't need to do it today，不能把it today 拆成一个短句，to do it应该作为一个短句。" +
		"8. 你需要确保每个短句都被拆出来，这很重要。" +
		"9. 注意拆分出单词时，不要一个单词拆分出多次。" +
		"10. 最后，输出完整的句子。" +
		"11. 输出时不要添加任何说明文字。" +
		"12. 输出为每个子句都有排序order、中文chinese、英文english、音标soundmark" +
		//"13. 区分哪组数据是哪个句子的，sentence键存句子，statements键存拆分的短句。" +
		//"14. 这种数据格式包含一个顶级的JSON对象，该对象包含一个名为sentences的键，其值是一个JSON数组。" +
		//"15. 这个JSON数组包含sentence、statements两个键，statements其值是一个JSON数组。" +
		//"16. 这个数组中的每个元素本身又是一个JSON数组。" +
		//"17. 每个子数组包含多个JSON对象，每个对象都有order、chinese、english和soundmark这些键。" +
		"请生成一个包含多个句子的对话结构数据，每个句子包含以下字段：" +
		"- sentence: 句子的英文文本" +
		"- statements: 一个包含多个语句的数组，每个语句包含以下字段：" +
		"  - chinese: 语句的中文翻译" +
		"  - english: 语句的英文文本" +
		"  - order: 语句在句子中的顺序" +
		"  - soundmark: 语句的音标" +
		"示例1：i like to play football on sunday with lucy." +
		"拆分为：I, like, to, I like to , play, football, play football, I like to play football, on, sunday , on Sunay," +
		" I like to play football on sunday, with, Lucy, with Lucy,  I like to play football on sunday with Lucy. " +
		"示例2:it is not important so i don't need to do it today." +
		"拆分为：it, is, it is, not, important, not important, it is not important, so ,I ,so I, don’t ,need ,don’t need, " +
		"to, do ,it ,to do it, today, it is not important so i don't need to do it today." +
		"如果你明白了，请将以下句子输出：" + sentencesStr

	//content = "You are a helpful assistant."

	// todo 对接deepseek 对话 API
	baseUrl := viper.GetString("deepSeekApi.baseUrl")
	apiKey := viper.GetString("deepSeekApi.apiKey")
	client := restclient.NewClient(baseUrl, 100)

	endpoint := "/v1/chat/completions"
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + apiKey
	headers["Content-Type"] = "application/json"

	// POST request example
	postPayload := DeepSeekRequest{
		Model: "deepseek-chat",
		Messages: []Message{
			{Role: "system", Content: content},
			//{Role: "user", Content: "Hello!"},
		},
		Stream: false,
	}
	//var postResult map[string]interface{}
	var deepSeekChatCompletionsResponse DeepSeekChatCompletionsResponse
	err := client.Post(endpoint, headers, postPayload, &deepSeekChatCompletionsResponse)
	if err != nil {
		fmt.Println("POST request failed:", err)
		return
	}
	fmt.Println("POST response:", deepSeekChatCompletionsResponse)

	responseContent := deepSeekChatCompletionsResponse.Choices[0].Message.Content*/

	responseContent := "```json\n{\n  \"sentence\": \"This is a very useful tool.\",\n  \"statements\": [\n    {\n      \"chinese\": \"这\",\n      \"english\": \"This\",\n      \"order\": 1,\n      \"soundmark\": \"ðɪs\"\n    },\n    {\n      \"chinese\": \"是\",\n      \"english\": \"is\",\n      \"order\": 2,\n      \"soundmark\": \"ɪz\"\n    },\n    {\n      \"chinese\": \"这是一个\",\n      \"english\": \"This is\",\n      \"order\": 3,\n      \"soundmark\": \"ðɪs ɪz\"\n    },\n    {\n      \"chinese\": \"一个\",\n      \"english\": \"a\",\n      \"order\": 4,\n      \"soundmark\": \"ə\"\n    },\n    {\n      \"chinese\": \"这是一个非常\",\n      \"english\": \"This is a\",\n      \"order\": 5,\n      \"soundmark\": \"ðɪs ɪz ə\"\n    },\n    {\n      \"chinese\": \"非常\",\n      \"english\": \"very\",\n      \"order\": 6,\n      \"soundmark\": \"ˈvɛri\"\n    },\n    {\n      \"chinese\": \"这是一个非常有用的\",\n      \"english\": \"This is a very\",\n      \"order\": 7,\n      \"soundmark\": \"ðɪs ɪz ə ˈvɛri\"\n    },\n    {\n      \"chinese\": \"有用的\",\n      \"english\": \"useful\",\n      \"order\": 8,\n      \"soundmark\": \"ˈjusfəl\"\n    },\n    {\n      \"chinese\": \"这是一个非常有用的工具\",\n      \"english\": \"This is a very useful\",\n      \"order\": 9,\n      \"soundmark\": \"ðɪs ɪz ə ˈvɛri ˈjusfəl\"\n    },\n    {\n      \"chinese\": \"工具\",\n      \"english\": \"tool\",\n      \"order\": 10,\n      \"soundmark\": \"tul\"\n    },\n    {\n      \"chinese\": \"这是一个非常有用的工具\",\n      \"english\": \"This is a very useful tool\",\n      \"order\": 11,\n      \"soundmark\": \"ðɪs ɪz ə ˈvɛri ˈjusfəl tul\"\n    }\n  ]\n}\n```"

	// 去掉 ```json 和 ``` 标记
	jsonStringContent := strings.TrimPrefix(responseContent, "```json")
	jsonStringContent = strings.TrimSuffix(jsonStringContent, "```")

	// 去掉换行符和多余的空格
	jsonStringContent = strings.ReplaceAll(jsonStringContent, "\n", "")
	//jsonString = strings.ReplaceAll(jsonString, " ", "")

	// 解析 JSON 数据
	var sentenceDatas []SentenceData
	if len(vo.StatementIds) == 0 {
		errSentenceDatas := json.Unmarshal([]byte(jsonStringContent), &sentenceDatas)
		if errSentenceDatas != nil {
			panic("Error unmarshaling JSON: %v" + errSentenceDatas.Error())
		}
	} else {
		var sentenceData SentenceData
		errSentenceDatas := json.Unmarshal([]byte(jsonStringContent), &sentenceData)
		if errSentenceDatas != nil {
			panic("Error unmarshaling JSON: %v" + errSentenceDatas.Error())
		}
		sentenceDatas = append(sentenceDatas, sentenceData)
	}

	// 输出结果
	/*for _, sentenceGroup := range response.Sentences {
		for _, sentence := range sentenceGroup {
			fmt.Printf("Order: %d, Chinese: %s, English: %s, Soundmark: %s\n", sentence.Order, sentence.Chinese, sentence.English, sentence.Soundmark)
		}
	}*/
	var sentence string
	var statements []Statement
	var statementIds []string
	for _, sentenceData := range sentenceDatas {
		sentence = strings.TrimSuffix(sentenceData.Sentence, ".")
		statements = sentenceData.Statements

		// todo 存入statements表
		var statementFirst earthworm.Statements
		sentenceFirst := dao.GetConn().Table("statements").
			Where("english = ?", sentence).
			First(&statementFirst)

		if sentenceFirst.Error != nil {
			if sentenceFirst.Error == gorm.ErrRecordNotFound {
				panic(base.ParamsError("sentence does not exist"))
			} else {
				panic("Error querying database: %v" + sentenceFirst.Error.Error())
			}
		} else {
			fmt.Println("sentence exists")

			// todo 插入短句
			var statementCreates []earthworm.Statements
			for _, statement := range statements {
				var statementCreate earthworm.Statements
				statementCreate.Id = base.GenerateUniqueTextID()
				statementCreate.Pid = statementFirst.Id
				statementCreate.CourseId = vo.CourseId
				statementCreate.Chinese = statement.Chinese
				statementCreate.English = statement.English
				statementCreate.Soundmark = statement.Soundmark
				statementCreate.Order = statement.Order
				statementCreates = append(statementCreates, statementCreate)
			}

			errCreate := dao.GetConn().Table("statements").
				Create(&statementCreates).Error
			if errCreate != nil {
				panic(base.ParamsErrorN())
			}
		}
	}

	/*var statementUpdates earthworm.Statements
	base.CopyProperties(&statementUpdates, vo)

	errUpdates := dao.GetConn().Table("statements").
		Updates(&statementUpdates).Error
	if errUpdates != nil {
		panic(base.ParamsErrorN())
	}*/

	c.JSON(http.StatusOK, base.Success(statementIds))
}
