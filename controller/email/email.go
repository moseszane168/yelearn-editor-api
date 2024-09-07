package email

import (
	"crf-mold/base"
	"crf-mold/common/constant"
	"crf-mold/dao"
	"crf-mold/model"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

const emailRe = `^[\.a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`

var re *regexp.Regexp

// @Tags 邮件
// @Summary 设置超时邮件发送人
// @Accept json
// @Produce json
// @Param Body body EmailSetVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {object} bool
// @Router /email/timeout [put]
func UpdateMoldTaskTimeoutEmail(c *gin.Context) {
	var vo EmailSetVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsErrorN())
	}

	t := vo.Type
	key := ""
	if t == "maintenance" {
		key = "maintenance_timeout_email"
	} else if t == "remodel" {
		key = "remodel_timeout_email"
	} else {
		panic(base.ResponseEnum[base.EMAIL_UNOKNOW_TYPE])
	}

	// 邮箱格式校验
	if vo.Receives != "" {
		emails := strings.Split(vo.Receives, ",")
		for _, v := range emails {
			if !validateEmailForamt(v) {
				panic(base.ResponseEnum[base.EMAIL_FORMAT_ERROR])
			}
		}
	}

	dao.GetConn().Table("properties").Where("`key` = ?", key).Updates(map[string]interface{}{
		"value":      vo.Receives,
		"updated_by": c.GetHeader(constant.USERID),
	})

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 邮件
// @Summary 获取超时邮件发送人
// @Accept json
// @Produce json
// @Param AuthToken header string false "Token"
// @Success 200 {object} EmailGetVO
// @Router /email/timeout [get]
func GetMoldTaskTimeoutEmail(c *gin.Context) {
	var mte model.Properties
	dao.GetConn().Table("properties").Where("`key` = 'maintenance_timeout_email' and is_deleted = 'N'").First(&mte)

	var rte model.Properties
	dao.GetConn().Table("properties").Where("`key` = 'remodel_timeout_email' and is_deleted = 'N'").First(&rte)

	var out EmailGetVO
	out.MaintenanceReceives = mte.Value
	out.RemodelReceives = rte.Value

	c.JSON(http.StatusOK, base.Success(out))
}

func validateEmailForamt(addr string) bool {
	if re == nil {
		re, _ = regexp.Compile(emailRe)
	}

	return re.Match([]byte(addr))
}
