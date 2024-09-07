/**
 * 用户中心
 */

package user

import (
	"crf-mold/base"
	"crf-mold/common"
	"crf-mold/common/constant"
	"crf-mold/dao"
	"crf-mold/model"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func getRsaPublicKey() string {
	return viper.GetString("secret.publicKey")
}

func getRsaPrivateKey() string {
	return viper.GetString("secret.privateKey")
}

// 解密密文
func DecrtptPassword(plaintext string) ([]byte, error) {
	publicKey := getRsaPublicKey()
	privateKey := getRsaPrivateKey()

	secretBytes, err := base64.StdEncoding.DecodeString(plaintext)
	if err != nil {
		panic(base.ResponseEnum[base.PARAMS_ERROR])
	}

	b, err := base.Decrypt([]byte(secretBytes), publicKey, privateKey)
	if err != nil {
		panic(base.ResponseEnum[base.PARAMS_ERROR])
	}

	return b, nil
}

// @Tags 登录
// @Summary 获取rsa公钥
// @Accept json
// @Produce json
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /user/rsa [get]
func GetRsaPublicKey(c *gin.Context) {
	c.JSON(http.StatusOK, base.Success(getRsaPublicKey()))
}

// @Tags 登录
// @Summary 密码加密
// @Accept json
// @Produce json
// @Param Body body EncryptPassWordVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /user/encrypt [post]
func EncryptPassword(c *gin.Context) {
	var vo EncryptPassWordVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	secrettext, err := base.Encrypt([]byte(vo.Plaintext), viper.GetString("secret.publicKey"))
	if err != nil {
		panic(err)
	}

	encodeString := base64.StdEncoding.EncodeToString(secrettext)

	c.JSON(http.StatusOK, encodeString)
}

// @Tags 登录
// @Summary 用户登录
// @Accept json
// @Produce json
// @Param Body body LoginVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /user/login [post]
func LoginIn(c *gin.Context) {
	// 参数处理
	var vo LoginVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	// 验证密码是否合法
	publicKey := getRsaPublicKey()
	privateKey := getRsaPrivateKey()
	if publicKey == "" || privateKey == "" {
		panic(base.UnknowError())
	}

	// 解密
	b, _ := DecrtptPassword(vo.PassWord)
	plaintext := string(b)

	// 验证用户是否存在
	md5Password := base.MD5(plaintext)
	var userInfo model.UserInfo
	if err := dao.GetConn().Table("user_info").Where("login_name = ? and is_deleted = 'N'", vo.LoginName).First(&userInfo).Error; err != nil {
		// 用户不存在
		panic(base.ResponseEnum[base.USER_NOT_EXIST])
	} else {
		// 用户已经锁住，无法登录
		if userInfo.IsLocked == "Y" {
			// 是否已经解锁
			if userInfo.UnlockTime.Before(time.Now()) {
				// 解锁
				UnlockUser(nil, userInfo.LoginName, true)
			} else {
				panic(base.ResponseEnum[base.USER_LOCKED])
			}
		}

		// 密码错误
		if userInfo.Password != md5Password {
			// 更新错误登录次数和锁定
			UpdateUserLockedInfo(vo.LoginName, c)
			panic(base.ResponseEnum[base.USER_OR_PASSWORD_ERROR])
		}
	}

	// 生成一个token
	token := base.MD5(base.RandStr(32))

	// 保存token
	TokenInfoMap.Store(token, &TokenInfo{
		LoginName: userInfo.LoginName,
		//ExpireTime: time.Now().Unix() + 7*24*3600, // 7天后过期
		ExpireTime: time.Now().Unix() + 2*3600, // 新需求：登录的token改成两小时
	})

	// 移除之前Token
	oldToken, ok := UserTokenMap.Load(vo.LoginName)
	if ok {
		TokenInfoMap.Delete(oldToken)
		logrus.Info("移除Token:", oldToken)
	}

	UserTokenMap.Store(vo.LoginName, token)

	// 登录成功，更新登录时间和次数
	if err := dao.GetConn().Table("login_info").Where("login_name = ? and is_deleted = 'N'", vo.LoginName).First(&model.LoginInfo{}).Error; err != nil {
		// 不存在，新增
		dao.GetConn().Create(&model.LoginInfo{
			LoginName:  vo.LoginName,
			FaultCount: 0,
		})
	} else {
		// 存在，更新
		dao.GetConn().Where("login_name = ?", vo.LoginName).Updates(&model.LoginInfo{
			FaultCount: 0,
		})
	}

	c.JSON(http.StatusOK, base.Success(map[string]string{
		"token": token,
	}))
}

// @Tags 登录
// @Summary 登出当前用户
// @Accept json
// @Produce json
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /user/logout [post]
func LogOut(c *gin.Context) {
	userId := c.Request.Header.Get(constant.USERID)
	token := c.Request.Header.Get("AuthToken")
	_, ok := TokenInfoMap.Load(token)
	if ok {
		TokenInfoMap.Delete(token)
		UserTokenMap.Delete(userId)
	}

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 用户管理
// @Summary 添加用户
// @Accept json
// @Produce json
// @Param Body body CreateUserVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /user [post]
func CreateUser(c *gin.Context) {
	var vo CreateUserVO

	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	// 登录名唯一
	var count int64
	dao.GetConn().Table("user_info").Where("login_name = ? and is_deleted = 'N'", vo.LoginName).Count(&count)
	if count > 0 {
		panic(base.ResponseEnum[base.USER_ID_REPEAT])
	}

	// 解密
	b, err := DecrtptPassword(vo.Password)
	if err != nil {
		panic(base.ParamsErrorN())
	}

	plaintext := string(b)
	if !ValidatePassWord(plaintext) {
		panic(base.ResponseEnum[base.PASSWORD_FORMAT_ERROR])
	}
	md5 := base.MD5(plaintext)

	var userinfo model.UserInfo
	base.CopyProperties(&userinfo, vo)

	userId := c.GetHeader(constant.USERID)
	userinfo.CreatedBy = userId
	userinfo.UpdatedBy = userId
	userinfo.Password = md5

	tx := dao.GetConn()
	defer dao.TransactionRollback(tx)

	tx.Table("user_info").Create(&userinfo)

	// 给新用户所有前端可配置的权限
	var authorities []model.UserAuthority
	tx.Table("user_authority").Where("is_deleted = 'N'").Where("display = 'Y'").Find(&authorities)
	if len(authorities) > 0 {
		userAuthorityRelList := make([]model.UserAuthorityRel, len(authorities))
		for i := 0; i < len(authorities); i++ {
			userAuthorityRelList[i] = model.UserAuthorityRel{
				LoginName:     userinfo.LoginName,
				AuthorityCode: authorities[i].Code,
				CreatedBy:     userId,
				UpdatedBy:     userId,
			}
		}
		tx.Table("user_authority_rel").CreateInBatches(userAuthorityRelList, len(userAuthorityRelList))
	}

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(userinfo.ID))
}

// @Tags 用户管理
// @Summary 更新用户
// @Accept json
// @Produce json
// @Param Body body UpdateUserVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /user [put]
func UpdateUser(c *gin.Context) {
	var vo UpdateUserVO

	if err := c.ShouldBindWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	// loginName不能重复
	var count int64
	dao.GetConn().Table("user_info").Where("id != ? and is_deleted = 'N' and login_name = ?", vo.ID, vo.LoginName).Count(&count)
	if count > 0 {
		panic(base.ResponseEnum[base.USER_ID_REPEAT])
	}

	dao.GetConn().Table("user_info").Where("id = ? and is_deleted = 'N'", vo.ID).Updates(&model.UserInfo{
		LoginName:  vo.LoginName,
		Name:       vo.Name,
		Department: vo.Department,
		UpdatedBy:  c.GetHeader(constant.USERID),
	})

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 用户管理
// @Summary 用户列表
// @Accept json
// @Produce json
// @Param name query string false "姓名或者工号"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /user [get]
func ListUser(c *gin.Context) {
	name := c.Query("name")

	var userList []model.UserInfo
	tx := dao.GetConn().Table("user_info").Where("is_deleted = 'N'")
	if name != "" {
		tx = tx.Where("(login_name like concat('%',?,'%') or name like concat('%',?,'%'))", name, name)
	}
	tx.Order("gmt_created desc").Find(&userList)

	if len(userList) > 0 {
		c.JSON(http.StatusOK, base.Success(userList))
	} else {
		c.JSON(http.StatusOK, base.Success([]model.UserInfo{}))
	}
}

// @Tags 用户管理
// @Summary 查询单个用户信息,管理员only
// @Accept json
// @Produce json
// @Param loginName query string false "loginName"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /user/one [get]
func SingelUserAdmin(c *gin.Context) {
	name := c.Query("loginName")
	if name == "" {
		panic(base.ParamsErrorN())
	}

	var result model.UserInfo
	if err := dao.GetConn().Table("user_info").Where("is_deleted = 'N'").Where("login_name like concat('%',?,'%')", name).First(&result).Error; err != nil {
		panic(base.ParamsErrorN())
	}

	c.JSON(http.StatusOK, base.Success(result))
}

// @Tags 用户管理
// @Summary 查询当前登录用户信息
// @Accept json
// @Produce json
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /user/self [get]
func LoginUserInfo(c *gin.Context) {
	userId := c.GetHeader(constant.USERID)

	var result LoginUserInfoVO
	var authority []model.UserAuthorityRel
	dao.GetConn().Table("user_info").Where("is_deleted = 'N' and login_name = ?", userId).Scan(&result)
	dao.GetConn().Table("user_authority_rel").Where("is_deleted = 'N' and login_name = ?", userId).Scan(&authority)

	var roles []string
	for _, v := range authority {
		roles = append(roles, v.AuthorityCode)
	}

	if len(roles) == 0 {
		result.Roles = nil
	} else {
		result.Roles = roles
	}

	c.JSON(http.StatusOK, base.Success(result))
}

// @Tags 用户管理
// @Summary 用户分页
// @Accept json
// @Produce json
// @Param name query string false "姓名或者工号"
// @Param currentPage query int false "当前页"
// @Param size query int false "每页数量"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /user/page [get]
func PageUser(c *gin.Context) {
	name := c.Query("name")

	current, size := base.GetPageParams(c)

	var userList []model.UserInfo
	tx := dao.GetConn().Table("user_info").Where("is_deleted = 'N'").Order("gmt_created desc")
	if name != "" {
		tx = tx.Where("(login_name like concat('%',?,'%') or name like concat('%',?,'%'))", name, name)
	}

	page := base.Page(tx, &userList, current, size)

	if len(userList) == 0 {
		page.List = []interface{}{}
	}

	c.JSON(http.StatusOK, base.Success(page))
}

// @Tags 用户管理
// @Summary 删除用户
// @Accept json
// @Produce json
// @Param Body body common.DeleteVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /user [delete]
func DeleteUser(c *gin.Context) {
	var vo common.DeleteVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	ids := vo.IDS
	for i := 0; i < len(ids); i++ {
		// 不能删除管理员
		if ids[i] == 1 {
			panic(base.ResponseEnum[base.ADMIN_CAN_NOT_DELETE])
		}
	}

	dao.GetConn().Table("user_info").Where("id in ?", ids).Updates(&model.SpareInfo{
		IsDeleted: "Y",
		UpdatedBy: c.GetHeader(constant.USERID),
	})

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 用户管理
// @Summary 重置密码,only管理员
// @Accept json
// @Produce json
// @Param Body body LoginVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /user/reset [put]
func ResetPassWord(c *gin.Context) {
	var vo LoginVO
	if err := c.ShouldBindWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	// 用户不存在
	var user model.UserInfo
	if err := dao.GetConn().Table("user_info").Where("login_name = ? and is_deleted = 'N'", vo.LoginName).First(&user).Error; err != nil {
		panic(base.ParamsErrorN())
	}

	// 解密
	b, err := DecrtptPassword(vo.PassWord)
	if err != nil {
		panic(base.ParamsErrorN())
	}

	// 密码格式
	plaintext := string(b)
	if !ValidatePassWord(plaintext) {
		panic(base.ResponseEnum[base.PASSWORD_FORMAT_ERROR])
	}

	// 更新密码
	md5Pwd := base.MD5(plaintext)
	dao.GetConn().Table("user_info").Where("id = ?", user.ID).Updates(&model.UserInfo{
		UpdatedBy: c.GetHeader(constant.USERID),
		Password:  md5Pwd,
		IsLocked:  "N",
	})

	// 更新登录信息
	dao.GetConn().Table("login_info").Where("login_name = ? and is_deleted = 'N'", vo.LoginName).Updates(&model.LoginInfo{
		FaultCount: 0,
	})

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 用户管理
// @Summary 修改当前用户密码
// @Accept json
// @Produce json
// @Param Body body ModifyPasswordVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /user/pwd [post]
func ModifyPassword(c *gin.Context) {
	var vo ModifyPasswordVO
	if err := c.ShouldBindWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	// 用户不存在
	var user model.UserInfo
	if err := dao.GetConn().Table("user_info").Where("login_name = ? and is_deleted = 'N'", c.GetHeader(constant.USERID)).First(&user).Error; err != nil {
		panic(base.ParamsErrorN())
	}

	// 解密旧密码
	old, err := DecrtptPassword(vo.OldPwd)
	if err != nil {
		panic(base.ParamsErrorN())
	}

	// 旧密码是否正确
	md5Pwd := base.MD5(string(old))
	if user.Password != md5Pwd {
		panic(base.ResponseEnum[base.USER_OR_PASSWORD_ERROR])
	}

	// 解密新密码
	new, err := DecrtptPassword(vo.NewPwd)
	if err != nil {
		panic(base.ParamsErrorN())
	}

	// 密码格式校验
	plaintext := string(new)
	if !ValidatePassWord(plaintext) {
		panic(base.ResponseEnum[base.PASSWORD_FORMAT_ERROR])
	}

	md5Pwd = base.MD5(plaintext)
	dao.GetConn().Table("user_info").Where("id = ?", user.ID).Updates(&model.UserInfo{
		UpdatedBy: c.GetHeader(constant.USERID),
		Password:  md5Pwd,
	})

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 权限
// @Summary 获取当前登录用户权限编码
// @Accept json
// @Produce json
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /user/authority [get]
func GetUserAuthority(c *gin.Context) {
	userId := c.GetHeader(constant.USERID)

	var result []model.UserAuthorityRel
	dao.GetConn().Table("user_authority_rel").Where("is_deleted = 'N' and login_name = ?", userId).Find(&result)

	var res []string
	for i := 0; i < len(result); i++ {
		res = append(res, result[i].AuthorityCode)
	}

	if len(res) > 0 {
		c.JSON(http.StatusOK, base.Success(res))
	} else {
		c.JSON(http.StatusOK, base.Success([]string{}))
	}
}

// @Tags 权限
// @Summary 获取指定用户权限编码,管理员only
// @Accept json
// @Produce json
// @Param loginName query string true "用户登录名"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /user/authority/one [get]
func GetUserAuthorityById(c *gin.Context) {
	userId := c.Query("loginName")
	if userId == "" {
		panic(base.ParamsErrorN())
	}

	var result []model.UserAuthorityRel
	dao.GetConn().Table("user_authority_rel").Where("is_deleted = 'N' and login_name = ?", userId).Find(&result)

	var res []string
	for i := 0; i < len(result); i++ {
		res = append(res, result[i].AuthorityCode)
	}

	if len(res) > 0 {
		c.JSON(http.StatusOK, base.Success(res))
	} else {
		c.JSON(http.StatusOK, base.Success([]string{}))
	}
}

// @Tags 权限
// @Summary 获取所有可配置的权限
// @Accept json
// @Produce json
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /user/authorities [get]
func GetAuthoritys(c *gin.Context) {
	var results []model.UserAuthority
	dao.GetConn().Table("user_authority").Where("is_deleted = 'N' and display = 'Y'").Order("code ASC").Find(&results)

	// 封装一下返回结果的格式
	res := make(map[string][]GetAuthoritysVO)
	for i := 0; i < len(results); i++ {
		list, ok := res[results[i].GroupName]
		if !ok {
			list = []GetAuthoritysVO{}
		}

		list = append(list, GetAuthoritysVO{
			Code:      results[i].Code,
			Name:      results[i].Name,
			GroupName: results[i].GroupName,
		})

		res[results[i].GroupName] = list
	}

	c.JSON(http.StatusOK, base.Success(res))
}

// @Tags 权限
// @Summary 更新指定用户权限
// @Accept json
// @Produce json
// @Param Body body UpdateUserAuthorityVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /user/authority [put]
func UpdateUserAuthority(c *gin.Context) {
	var vo UpdateUserAuthorityVO
	if err := c.ShouldBindWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	// 用户不存在
	var count int64
	dao.GetConn().Table("user_info").Where("is_deleted = 'N' and login_name = ?", vo.LoginName).Count(&count)
	if count == 0 {
		panic(base.ParamsErrorN())
	}

	// 只能配置可配置的权限
	//dao.GetConn().Table("user_authority").Where("is_deleted = 'N' and display = 'Y'").Where("code in ?", vo.Codes).Count(&count)
	//if count != int64(len(vo.Codes)) {
	//	panic(base.ParamsErrorN())
	//}

	// 手动开启事务和指定事务归滚方法，后续代码的panic将触发事务回滚方法TransactionRollback，抛出999异常
	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	// 先删后增
	if err := tx.Table("user_authority_rel").Where("login_name = ?", vo.LoginName).Updates(&model.UserAuthorityRel{
		UpdatedBy: c.GetHeader(constant.USERID),
		IsDeleted: "Y",
	}).Error; err != nil {
		panic(err)
	}

	userId := c.GetHeader(constant.USERID)

	// 封装entity
	if len(vo.Codes) > 0 {
		var list []model.UserAuthorityRel
		for i := 0; i < len(vo.Codes); i++ {
			code := vo.Codes[i]
			list = append(list, model.UserAuthorityRel{
				LoginName:     vo.LoginName,
				AuthorityCode: code,
				CreatedBy:     userId,
				UpdatedBy:     userId,
			})
		}

		// 批量保存
		if err := tx.Table("user_authority_rel").CreateInBatches(list, len(list)).Error; err != nil {
			panic(err)
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(true))
}
