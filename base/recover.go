/**
 * 统一异常处理
 */

package base

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Recover(c *gin.Context) {
	defer func() {
		// 取执行的错误信息
		if r := recover(); r != nil {
			// //打印错误堆栈信息
			debug.PrintStack()
			logrus.WithField("stack", string(debug.Stack())).Errorf("panic: %v\n", r)

			response, ok := r.(*Response)

			if !ok {
				c.JSON(http.StatusInternalServerError, UnknowError())
			} else {
				c.JSON(http.StatusOK, response)
			}

			//终止后续接口调用，不加的话recover到异常后，还会继续执行接口里后续代码
			c.Abort()
		}
	}()
	//加载完 defer recover，继续后续接口调用
	c.Next()
}
