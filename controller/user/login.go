package user

import (
	"crf-mold/base"
	"crf-mold/dao"
	"crf-mold/model"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// 密码格式校验
func ValidatePassWord(pwd string) bool {
	// 规则:数字字母中英文符号组成,6-12位
	res, err := regexp.Match("^[a-zA-Z0-9`\\\\~!@#$%^&*()_\\-+=<>?:\"{}|,.\\/;'\\[\\]·~！@#￥%……&*（）——\\-+={}|《》？：“”【】、；‘'，。、]{6,12}$", []byte(pwd))
	if err != nil {
		logrus.WithError(err).Info("ValidatePassWord error")
		return false
	}

	return res
}

// 更新登录失败次数
func UpdateLoginFaultCount(tx *gorm.DB, loginName string, faultCount int) {
	if tx == nil {
		tx = dao.GetConn()
	}
	tx.Table("login_info").Where("login_name = ?", loginName).Updates(map[string]interface{}{
		"updated_by":  loginName,
		"fault_count": faultCount,
	})
}

// 更新用户锁定信息
func UpdateUserLockedInfo(loginName string, c *gin.Context) {
	// 用户不存在
	var userInfo model.UserInfo
	if err := dao.GetConn().Table("user_info").Where("login_name = ? and is_deleted = 'N'", loginName).First(&userInfo).Error; err != nil {
		return
	}

	var loginInfo model.LoginInfo
	// 更新登录时间和次数
	if err := dao.GetConn().Table("login_info").Where("login_name = ? and is_deleted = 'N'", loginName).First(&loginInfo).Error; err != nil {
		// 不存在，新增
		dao.GetConn().Table("login_info").Create(&model.LoginInfo{
			LoginName:  loginName,
			FaultCount: 1,
			CreatedBy:  "system",
			UpdatedBy:  "system",
		})
	} else {
		// 存在，更新
		count := loginInfo.FaultCount

		// 是否今天连错
		if base.SameDayWithNow(loginInfo.LastLoginDate) {
			count++
		} else {
			count = 1
		}

		// 登录失败次数超过三次，进行锁定
		if count > 3 {
			UnlockUser(nil, loginName, false)
			// 更新登录错误次数
			UpdateLoginFaultCount(nil, loginName, 0)
			panic(base.ResponseEnum[base.USER_LOCKED])
		} else {
			// 更新登录错误次数
			UpdateLoginFaultCount(nil, loginName, count)
		}
	}
}

// 解锁和锁定
func UnlockUser(tx *gorm.DB, loginName string, unlock bool) {
	if tx == nil {
		tx = dao.GetConn()
	}

	if unlock { // 解锁用户
		tx.Table("user_info").Where("login_name = ?", loginName).Updates(&model.UserInfo{
			IsLocked: "N",
		})
	} else { // 锁定用户10分钟
		unlockTime := time.Now().Add(time.Minute * 10)
		tx.Table("user_info").Where("login_name = ?", loginName).Updates(&model.UserInfo{
			IsLocked:   "Y",
			UnlockTime: &unlockTime,
		})
	}
}
