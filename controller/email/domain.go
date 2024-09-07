package email

import (
	"bytes"
	"crf-mold/base"
	"crf-mold/dao"
	"crf-mold/model"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"

	"github.com/jordan-wright/email"
	"github.com/spf13/viper"
)

type TIMEOUT_TYPE string

const (
	REMODEL_TIMEOUT     TIMEOUT_TYPE = "remodel_timeout"
	MAINTENANCE_TIMEOUT TIMEOUT_TYPE = "maintenance_timeout"
	REMODEL_CREATE      TIMEOUT_TYPE = "remodel_create"
	MAINTENANCE_CREATE  TIMEOUT_TYPE = "maintenance_create"
	EmailTemplate       string       = `<!DOCTYPE html>
					<html>
						<head>
							<style>
								table {
								  font-family: arial, sans-serif;
								  border-collapse: collapse;
                                  font-size: 10px;
								  width: 100%;
								}
								td, th {
								  border: 1px solid #dddddd;
								  text-align: left;
								  padding: 8px;
								}
								.button {
									display: block;
									width: 150px;
									height: 25px;
									background: #4E9CAF;
									padding: 10px;
									text-align: center;
									border-radius: 5px;
									color: white;
									font-weight: bold;
									line-height: 25px;
								}
								tr:nth-child(even) {
								  background-color: #dddddd;
								}
							</style>
						</head>
						<body>
							{{ .tableName }}
							<table>
								<tr>
									{{ .tableHead }}
								</tr>
									{{ .tableData }}
							</table>
							{{ .link }}
						</body>
					</html>
`
)

type Email struct {
	Title        string
	TableName    string
	TableHeaders string
	Content      string
	EmailContent string
	Link         string
	UseHtml      bool
}

var emailFormatMap = map[TIMEOUT_TYPE]Email{
	REMODEL_TIMEOUT: {
		Title:        "模具管理改造任务超时提醒",
		TableName:    "模具管理已超时改造任务",
		TableHeaders: `<th>序号</th><th>模具编号</th><th>零件号</th><th>项目名称</th><th>改造开始时间</th><th>改造结束时间</th><th>责任人</th>`,
		Link:         `</br><a href="http://10.102.1.174/transformation" class="button">前往模具管理系统</a>`,
		UseHtml:      true,
		Content:      "模具超期改造任务超时，改造编号:%s",
	},
	MAINTENANCE_TIMEOUT: {
		Title:        "模具管理保养任务超时提醒",
		TableName:    "模具管理已超时保养任务",
		TableHeaders: `<th>序号</th><th>模具编号</th><th>零件号</th><th>项目名称</th><th>创建时间</th><th>保养人</th>`,
		Link:         `</br><a href="http://10.102.1.174/moldMaintain" class="button">前往模具管理系统</a>`,
		UseHtml:      true,
		Content:      "有一条保养任务超时，保养编号：%s",
	},
	REMODEL_CREATE: {
		Title:        "模具管理改造任务创建提醒",
		TableName:    "模具管理待改造任务",
		TableHeaders: `<th>序号</th><th>模具编号</th><th>零件号</th><th>项目名称</th><th>改造开始时间</th><th>改造结束时间</th><th>责任人</th>`,
		Link:         `</br><a href="http://10.102.1.174/transformation" class="button">前往模具管理系统</a>`,
		UseHtml:      true,
		Content:      "有一条改造任务待完成，改造编号:%s",
	},
	MAINTENANCE_CREATE: {
		Title:        "模具管理保养任务创建提醒",
		TableName:    "模具管理待保养任务",
		TableHeaders: `<th>序号</th><th>模具编号</th><th>零件号</th><th>项目名称</th><th>创建时间</th><th>保养人</th>`,
		Link:         `</br><a href="http://10.102.1.174/moldMaintain" class="button">前往模具管理系统</a>`,
		UseHtml:      true,
		Content:      "有一条保养任务待完成，保养编号：%s",
	},
}

func FormatMessage(typ TIMEOUT_TYPE, argument ...interface{}) Email {
	title := emailFormatMap[typ].Title
	content := fmt.Sprintf(emailFormatMap[typ].Content, argument...)
	return Email{
		Title:   title,
		Content: content,
	}
}

func formatMaintenanceMessage(tasks interface{}) string {
	var dataString string
	fmtString := "<tr><td>%s</td>\n<td>%s</td>\n<td>%s</td>\n<td>%s</td>\n<td>%s</td>\n<td>%s</td></tr>\n"
	if tasks, ok := tasks.([]MaintenanceTaskEmailVO); ok {
		for _, task := range tasks {
			dataString += fmt.Sprintf(fmtString, task.TaskCode, task.MoldCode, task.PartCodes, task.ProjectName, task.GmtCreated, task.Operator)
		}
	}

	return dataString
}

func formatRemodelMessage(tasks interface{}) string {
	var dataString string
	fmtString := "<tr><td>%s</td>\n<td>%s</td>\n<td>%s</td>\n<td>%s</td>\n<td>%s</td>\n<td>%s</td>\n<td>%s</td></tr>\n"
	if tasks, ok := tasks.([]RemodelEmailVO); ok {
		for _, task := range tasks {
			dataString += fmt.Sprintf(fmtString, task.Code, task.MoldCode, task.PartCodes, task.ProjectName, task.RemodelStartTime, task.RemodelEndTime, task.Director)
		}
	}
	return dataString
}

func FormatEmailMessage(typ TIMEOUT_TYPE, result interface{}) Email {
	var dataString string
	title := emailFormatMap[typ].Title
	buf := new(bytes.Buffer)
	tmp, _ := template.New("emailTemplate").Parse(EmailTemplate)

	if base.StrIn(string(typ), []string{string(MAINTENANCE_CREATE), string(MAINTENANCE_TIMEOUT)}) {
		dataString = formatMaintenanceMessage(result)
	} else {
		dataString = formatRemodelMessage(result)
	}

	tmp.Execute(buf, map[string]interface{}{
		"tableName": emailFormatMap[typ].TableName,
		"tableHead": emailFormatMap[typ].TableHeaders,
		"link":      emailFormatMap[typ].Link,
		"tableData": dataString,
	})

	return Email{
		Title:        title,
		EmailContent: buf.String(),
		UseHtml:      emailFormatMap[typ].UseHtml,
	}
}

func SendEmail(address string, typ TIMEOUT_TYPE, result interface{}) error {
	em := email.NewEmail()
	// 设置 sender 发送方 的邮箱 ， 此处可以填写自己的邮箱
	em.From = viper.GetString("email.userName")

	// 设置 receiver 接收方 的邮箱  此处也可以填写自己的邮箱， 就是自己发邮件给自己
	em.To = []string{address}

	msg := FormatEmailMessage(typ, result)

	// 设置主题
	em.Subject = msg.Title

	if msg.UseHtml {
		em.HTML = []byte(msg.EmailContent)
	} else {
		em.Text = []byte(msg.Content)
	}

	//设置服务器相关的配置
	addr := fmt.Sprintf("%s:%d", viper.GetString("email.host"), viper.GetInt("email.port"))

	// 没配账号密码则匿名无授权发送
	var auth smtp.Auth
	if viper.GetString("email.passWord") != "" {
		auth = smtp.PlainAuth("", viper.GetString("email.userName"), viper.GetString("email.passWord"), viper.GetString("email.host"))
	}

	err := em.Send(addr, auth)
	if err != nil {
		return err
	}

	return nil
}

func GetMoldTaskTimeoutEmailReceiver(typ TIMEOUT_TYPE) ([]string, error) {
	t := typ
	key := ""
	if t == MAINTENANCE_TIMEOUT {
		key = "maintenance_timeout_email"
	} else if t == REMODEL_TIMEOUT {
		key = "remodel_timeout_email"
	} else {
		panic(base.ResponseEnum[base.EMAIL_UNOKNOW_TYPE])
	}

	var enitty model.Properties
	if err := dao.GetConn().Table("properties").Where("`key` = ?", key).First(&enitty).Error; err != nil {
		return nil, err
	}

	return strings.Split(enitty.Value, ","), nil
}
