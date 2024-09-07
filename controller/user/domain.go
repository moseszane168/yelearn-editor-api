package user

import (
	"crf-mold/dao"
	"crf-mold/model"
)

// 获取所有用户名称
func GetAllUserLoginName() []string {
	var result []model.UserInfo
	dao.GetConn().Table("user_info").Where("is_deleted = 'N'").Find(&result)

	res := make([]string, len(result))
	for i := 0; i < len(result); i++ {
		res[i] = result[i].LoginName
	}

	return res
}
