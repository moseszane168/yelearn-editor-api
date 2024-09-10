package statements

import (
	"bytes"
	"crf-mold/base"
	restclient "crf-mold/common/http"
	"crf-mold/dao"
	"crf-mold/model/earthworm"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/spf13/viper"
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
	id := c.Query("id")
	if id == "" {
		panic(base.ParamsErrorN())
	}

	//userId := c.GetHeader(constant.USERID)

	var statement earthworm.Statements
	errDelete := dao.GetConn().Table("statements").
		Where("id = ?", id).
		Delete(&statement).Error
	if errDelete != nil {
		panic(base.ParamsError(errDelete.Error()))
	}

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

	var statement []earthworm.Statements
	errFind := dao.GetConn().Table("statements").
		Where("course_id = ?", vo.CourseId).
		Where("pid = ?", "").
		//Order("`order` desc").
		Find(&statement).Error
	if errFind != nil {
		panic(base.ParamsError(errFind.Error()))
	}

	if len(statement) > 0 {
		c.JSON(http.StatusOK, base.Success(statement))
	} else {
		c.JSON(http.StatusOK, base.Success([]earthworm.Courses{}))
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
		"如果你明白了，请将以下句子输出：" + vo.Statement

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
			{Role: "user", Content: "Hello!"},
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

	responseContent := deepSeekChatCompletionsResponse.Choices[0].Message.Content

	//responseString := "```json\n[\n    {\n        \"sentence\": \"it is not important so i don't need to do it today.\",\n        \"statements\": [\n            {\n                \"chinese\": \"它\",\n                \"english\": \"it\",\n                \"order\": 1,\n                \"soundmark\": \"/ɪt/\"\n            },\n            {\n                \"chinese\": \"是\",\n                \"english\": \"is\",\n                \"order\": 2,\n                \"soundmark\": \"/ɪz/\"\n            },\n            {\n                \"chinese\": \"它是不\",\n                \"english\": \"it is\",\n                \"order\": 3,\n                \"soundmark\": \"/ɪt ɪz/\"\n            },\n            {\n                \"chinese\": \"不\",\n                \"english\": \"not\",\n                \"order\": 4,\n                \"soundmark\": \"/nɒt/\"\n            },\n            {\n                \"chinese\": \"重要\",\n                \"english\": \"important\",\n                \"order\": 5,\n                \"soundmark\": \"/ɪmˈpɔːrtnt/\"\n            },\n            {\n                \"chinese\": \"不重要\",\n                \"english\": \"not important\",\n                \"order\": 6,\n                \"soundmark\": \"/nɒt ɪmˈpɔːrtnt/\"\n            },\n            {\n                \"chinese\": \"它是不重要的\",\n                \"english\": \"it is not important\",\n                \"order\": 7,\n                \"soundmark\": \"/ɪt ɪz nɒt ɪmˈpɔːrtnt/\"\n            },\n            {\n                \"chinese\": \"所以\",\n                \"english\": \"so\",\n                \"order\": 8,\n                \"soundmark\": \"/səʊ/\"\n            },\n            {\n                \"chinese\": \"我\",\n                \"english\": \"i\",\n                \"order\": 9,\n                \"soundmark\": \"/aɪ/\"\n            },\n            {\n                \"chinese\": \"所以我\",\n                \"english\": \"so i\",\n                \"order\": 10,\n                \"soundmark\": \"/səʊ aɪ/\"\n            },\n            {\n                \"chinese\": \"不需要\",\n                \"english\": \"don’t need\",\n                \"order\": 11,\n                \"soundmark\": \"/dəʊnt niːd/\"\n            },\n            {\n                \"chinese\": \"做\",\n                \"english\": \"to do\",\n                \"order\": 12,\n                \"soundmark\": \"/tuː duː/\"\n            },\n            {\n                \"chinese\": \"它\",\n                \"english\": \"it\",\n                \"order\": 13,\n                \"soundmark\": \"/ɪt/\"\n            },\n            {\n                \"chinese\": \"今天\",\n                \"english\": \"today\",\n                \"order\": 14,\n                \"soundmark\": \"/təˈdeɪ/\"\n            },\n            {\n                \"chinese\": \"它是不重要的所以我今天不需要做它\",\n                \"english\": \"it is not important so i don't need to do it today\",\n                \"order\": 15,\n                \"soundmark\": \"/ɪt ɪz nɒt ɪmˈpɔːrtnt səʊ aɪ dəʊnt niːd tuː duː ɪt təˈdeɪ/\"\n            }\n        ]\n    },\n    {\n        \"sentence\": \"i like to play football on sunday with lucy.\",\n        \"statements\": [\n            {\n                \"chinese\": \"我\",\n                \"english\": \"i\",\n                \"order\": 1,\n                \"soundmark\": \"/aɪ/\"\n            },\n            {\n                \"chinese\": \"喜欢\",\n                \"english\": \"like\",\n                \"order\": 2,\n                \"soundmark\": \"/laɪk/\"\n            },\n            {\n                \"chinese\": \"我喜欢\",\n                \"english\": \"i like\",\n                \"order\": 3,\n                \"soundmark\": \"/aɪ laɪk/\"\n            },\n            {\n                \"chinese\": \"去\",\n                \"english\": \"to\",\n                \"order\": 4,\n                \"soundmark\": \"/tuː/\"\n            },\n            {\n                \"chinese\": \"我喜欢去\",\n                \"english\": \"i like to\",\n                \"order\": 5,\n                \"soundmark\": \"/aɪ laɪk tuː/\"\n            },\n            {\n                \"chinese\": \"玩\",\n                \"english\": \"play\",\n                \"order\": 6,\n                \"soundmark\": \"/pleɪ/\"\n            },\n            {\n                \"chinese\": \"足球\",\n                \"english\": \"football\",\n                \"order\": 7,\n                \"soundmark\": \"/ˈfʊtbɔːl/\"\n            },\n            {\n                \"chinese\": \"玩足球\",\n                \"english\": \"play football\",\n                \"order\": 8,\n                \"soundmark\": \"/pleɪ ˈfʊtbɔːl/\"\n            },\n            {\n                \"chinese\": \"我喜欢去玩足球\",\n                \"english\": \"i like to play football\",\n                \"order\": 9,\n                \"soundmark\": \"/aɪ laɪk tuː pleɪ ˈfʊtbɔːl/\"\n            },\n            {\n                \"chinese\": \"在\",\n                \"english\": \"on\",\n                \"order\": 10,\n                \"soundmark\": \"/ɒn/\"\n            },\n            {\n                \"chinese\": \"星期天\",\n                \"english\": \"sunday\",\n                \"order\": 11,\n                \"soundmark\": \"/ˈsʌndeɪ/\"\n            },\n            {\n                \"chinese\": \"在星期天\",\n                \"english\": \"on sunday\",\n                \"order\": 12,\n                \"soundmark\": \"/ɒn ˈsʌndeɪ/\"\n            },\n            {\n                \"chinese\": \"我喜欢去玩足球在星期天\",\n                \"english\": \"i like to play football on sunday\",\n                \"order\": 13,\n                \"soundmark\": \"/aɪ laɪk tuː pleɪ ˈfʊtbɔːl ɒn ˈsʌndeɪ/\"\n            },\n            {\n                \"chinese\": \"和\",\n                \"english\": \"with\",\n                \"order\": 14,\n                \"soundmark\": \"/wɪð/\"\n            },\n            {\n                \"chinese\": \"露西\",\n                \"english\": \"lucy\",\n                \"order\": 15,\n                \"soundmark\": \"/ˈluːsi/\"\n            },\n            {\n                \"chinese\": \"和露西\",\n                \"english\": \"with lucy\",\n                \"order\": 16,\n                \"soundmark\": \"/wɪð ˈluːsi/\"\n            },\n            {\n                \"chinese\": \"我喜欢去玩足球在星期天和露西\",\n                \"english\": \"i like to play football on sunday with lucy\",\n                \"order\": 17,\n                \"soundmark\": \"/aɪ laɪk tuː pleɪ ˈfʊtbɔːl ɒn ˈsʌndeɪ wɪð ˈluːsi/\"\n            }\n        ]\n    }\n]\n```"

	// 去掉 ```json 和 ``` 标记
	jsonStringContent := strings.TrimPrefix(responseContent, "```json")
	jsonStringContent = strings.TrimSuffix(jsonStringContent, "```")

	// 去掉换行符和多余的空格
	jsonStringContent = strings.ReplaceAll(jsonStringContent, "\n", "")
	//jsonString = strings.ReplaceAll(jsonString, " ", "")

	// 解析 JSON 数据
	var sentenceDatas []SentenceData
	errSentenceDatas := json.Unmarshal([]byte(jsonStringContent), &sentenceDatas)
	if errSentenceDatas != nil {
		panic("Error unmarshaling JSON: %v" + errSentenceDatas.Error())
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
			for _, statement := range statements {
				var statementCreate earthworm.Statements
				statementCreate.Id = base.GenerateUniqueTextID()
				statementCreate.Pid = statementFirst.Id
				statementCreate.CourseId = vo.CourseId
				statementCreate.Chinese = statement.Chinese
				statementCreate.English = statement.English
				statementCreate.Soundmark = statement.Soundmark
				statementCreate.Order = statement.Order
				errCreate := dao.GetConn().Table("statements").
					Create(&statementCreate).Error
				if errCreate != nil {
					panic(base.ParamsErrorN())
				}

				statementIds = append(statementIds, statementCreate.Id)
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
