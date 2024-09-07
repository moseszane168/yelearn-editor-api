package email_test

import (
	"log"
	"net/smtp"
	"testing"

	"github.com/jordan-wright/email"
)

func TestSendEmail(t *testing.T) {
	em := email.NewEmail()
	// 设置 sender 发送方 的邮箱 ， 此处可以填写自己的邮箱
	em.From = "chengong@crf-tech.com"

	// 设置 receiver 接收方 的邮箱  此处也可以填写自己的邮箱， 就是自己发邮件给自己
	em.To = []string{"980925441@qq.com"}

	// 设置主题
	em.Subject = "测试邮件发送"

	// 简单设置文件发送的内容，暂时设置成纯文本
	em.Text = []byte("模具管理任务超时提醒")

	//设置服务器相关的配置
	err := em.Send("imaphz.qiye.163.com:25", smtp.PlainAuth("", "chengong@crf-tech.com", "8fGEZmjaRxGtrsn4", "imaphz.qiye.163.com"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("send successfully ... ")
}
