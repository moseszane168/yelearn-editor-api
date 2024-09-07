/**
 * 用户权限
 */

package user

import (
	"crf-mold/base"
	"crf-mold/common/constant"
	"crf-mold/dao"
	"crf-mold/model"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthorityVO struct {
	Uri    string
	Method string
}

type TokenInfo struct {
	LoginName  string `json:"loginName"`
	ExpireTime int64  `json:"expireTime"`
}

// Token对应的用户ID
var TokenInfoMap sync.Map

// 用户ID对应的Token
var UserTokenMap sync.Map

// addFixToken 添加固定Token
func addFixToken() {
	token := "bbfdca13472f1189a596fe1ee856ec71"
	loginName := "appAccess"

	// 保存token
	TokenInfoMap.Store(token, &TokenInfo{
		LoginName:  loginName,
		ExpireTime: time.Now().Unix() + 24*3600*365*10, // 10年
	})
}

func Authority(c *gin.Context) {
	// TODO:这里不支持路径参数
	uri := c.Request.RequestURI
	uri = strings.Split(uri, "?")[0]

	allowedUriArray := []string{"/v1/user/rsa", "/v1/user/login", "/v1/template", "/v1/file", "/v1/system/ws", "/v1/dict/all", "/v1/dict", "/v1/mold"}
	// 放行特殊接口
	if base.StrIn(uri, allowedUriArray) || strings.HasPrefix(uri, "/v1/ws") ||
		strings.HasPrefix(uri, "/v1/rtsp/upload") || strings.HasPrefix(uri, "/v1/rtsp/live") {
		// 继续调用
		c.Next()
		return
	}

	token := c.Request.Header.Get("AuthToken")

	// 未传token,401未登录授权访问
	if token == "" {
		c.JSON(http.StatusUnauthorized, base.ResponseEnum[base.UNAUTHORIZED])
		c.Abort()
		return
	}

	v, ok := TokenInfoMap.Load(token)
	// token无效,401未登录授权访问
	if !ok {
		c.JSON(http.StatusUnauthorized, base.ResponseEnum[base.UNAUTHORIZED])
		c.Abort()
		return
	}

	tokenInfo := v.(*TokenInfo)

	// token是否过期
	now := time.Now().Unix()
	diff := tokenInfo.ExpireTime - now

	// 过期了
	if diff <= 0 {
		c.JSON(http.StatusUnauthorized, base.ResponseEnum[base.UNAUTHORIZED])
		c.Abort()
		return
	}

	// 有效期10分钟内,续过期时间到7天
	if diff <= 10*60 {
		tokenInfo.ExpireTime = time.Now().Unix() + 7*24*3600
	}

	userId := tokenInfo.LoginName
	method := c.Request.Method

	// 获得所有权限
	authoritys := GetAllAuthority()

	// 如果访问的接口被管控
	if _, ok := authoritys[uri+method]; ok {
		// 管理员和App访问直接放过
		if !base.StrIn(userId, []string{"admin", "appAccess"}) {
			// 非管理员判断是否有权限访问
			var authorityList []AuthorityVO
			dao.GetConn().Raw(`
				SELECT
					ua.uri,
					ua.method 
				FROM
					user_info ui
					INNER JOIN user_authority_rel uar ON uar.login_name = ui.login_name
					AND uar.is_deleted = 'N'
					INNER JOIN user_authority ua ON ua.code = uar.authority_code
					AND ua.is_deleted = 'N' 
				WHERE
					ui.is_deleted = 'N' and ua.uri = ? and ua.method = ? and ui.login_name = ?`, uri, method, userId).Scan(&authorityList)

			// 终止
			if len(authorityList) == 0 {
				c.JSON(http.StatusForbidden, base.ResponseEnum[base.NO_PERSSION])
				c.Abort()
				return
			}
		}
	}

	// 添加登录用户ID到header
	c.Request.Header.Set(constant.USERID, userId)

	// 继续调用
	c.Next()
}

func GetAllAuthority() map[string]model.UserAuthority {
	var authoritys []model.UserAuthority
	dao.GetConn().Table("user_authority").Where("is_deleted = 'N'").Find(&authoritys)

	resultMap := make(map[string]model.UserAuthority)

	for i := 0; i < len(authoritys); i++ {
		authority := authoritys[i]
		resultMap[authority.URI+authority.Method] = authority
	}

	return resultMap
}

// InitToken 初始化Token(定期清理过期Token, 添加固定Token)
func InitToken() {
	//添加固定Token
	addFixToken()

	go func() {
		defer func() {
			if err := recover(); err != nil {
				logrus.Error("InitToken error Panic info is: %v", err)
				debug.PrintStack()
				logrus.WithField("stack", string(debug.Stack())).Errorf("panic: %v\n", err)
			}
		}()

		for {
			TokenInfoMap.Range(func(key, value interface{}) bool {
				tokeninfo := value.(*TokenInfo)
				now := time.Now().Unix()
				if tokeninfo.ExpireTime < now {
					TokenInfoMap.Delete(key)
					logrus.Info("删除过期key:", key)
				}
				return true
			})
			// 十分钟查看一次，进行清理
			time.Sleep(10 * time.Minute)
		}
	}()
}
