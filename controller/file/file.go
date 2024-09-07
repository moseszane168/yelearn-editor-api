package file

import (
	"crf-mold/base"
	minioutil "crf-mold/common/minio"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// 随机生成文件名
func RandomFileName(fileName string) string {
	prefixIndex := strings.LastIndex(fileName, ".")
	// 后缀 .xxx
	lastfix := fileName[prefixIndex:]
	return base.UUID() + lastfix
}

// @Tags 文件
// @Summary 文件上传
// @Accept json
// @Produce json
// @Param file formData file true "文件"
// @Param AuthToken header string false "Token"
// @Success 200 {object} file.FileOutVO
// @Router /file [post]
func UploadFile(c *gin.Context) {
	header, err := c.FormFile("file")
	if err != nil {
		panic(base.ParamsErrorN())
	}

	// 随机文件名
	dst := RandomFileName(header.Filename)

	// 保存上传的文件到本地
	if err := c.SaveUploadedFile(header, dst); err != nil {
		panic(err)
	}

	// 结束后删除文件
	defer func() {
		err = os.Remove(dst)
		if err != nil {
			panic(err)
		}
	}()

	// 上传Minio
	info := minioutil.UploadMinio(dst, dst)

	// 返回Key
	c.JSON(http.StatusOK, base.Success(&FileOutVO{
		BucketName: info.Bucket,
		FileKey:    info.Key,
		FileName:   header.Filename,
	}))
}
