package user

import (
	"crf-mold/base"
	"time"
)

type CreateUserVO struct {
	LoginName  string `json:"loginName" binding:"required"`  // 工号
	Name       string `json:"name" binding:"required"`       // 姓名
	Department string `json:"department" binding:"required"` // 部门
	Password   string `json:"password" binding:"required"`   // 密码，rsa加密
} // @name CreateUserVO

type UpdateUserVO struct {
	ID         int    `json:"id" binding:"required"`         // ID
	LoginName  string `json:"loginName" binding:"required"`  // 工号
	Name       string `json:"name" binding:"required"`       // 姓名
	Department string `json:"department" binding:"required"` // 部门
} // @name UpdateUserVO

type LoginVO struct {
	LoginName string `json:"loginName" binding:"required"` // 工号
	PassWord  string `json:"passWord" binding:"required"`  // 密码，rsa加密
} // @name LoginVO

type UpdateUserAuthorityVO struct {
	LoginName string   `json:"loginName" binding:"required"` // 工号
	Codes     []string `json:"codes" binding:"min=1"`        // 权限编码
} // @name UpdateUserAuthorityVO

type GetAuthoritysVO struct {
	Code      string `json:"code"`      // 权限编码
	Name      string `json:"name"`      // 权限说明
	GroupName string `json:"groupName"` // 权限组
} // @name GetAuthoritysVO

type ModifyPasswordVO struct {
	OldPwd string `json:"oldPwd" binding:"required"` // 原始密码，rsa加密
	NewPwd string `json:"newPwd" binding:"required"` // 新密码，rsa加密
} // @name ModifyPasswordVO

type EncryptPassWordVO struct {
	Plaintext string `json:"plaintext" binding:"required"` // 明文密码
} // @name EncryptPassWordVO

type LoginUserInfoVO struct {
	ID         int64      `json:"id"`
	LoginName  string     `json:"loginName"`   // 工号
	Name       string     `json:"name"`        // 姓名
	Department string     `json:"department"`  // 部门
	IsRoot     string     `json:"isRoot"`      // 是否超级管理员
	IsLocked   string     `json:"isLocked"`    // 是否锁定
	IsDeleted  string     `json:"isDeleted"`   // 是否删除
	UnlockTime *time.Time `json:"unlock_time"` // 解锁时间
	GmtCreated base.Time  `json:"gmtCreated"`  // 创建时间
	GmtUpdated base.Time  `json:"gmtUpdated"`  // 修改时间
	Roles      []string   `json:"roles"`       // 权限
} // @name LoginUserInfoVO
