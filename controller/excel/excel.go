/**
 * Excel模块
 */

package excel

import (
	"crf-mold/base"
	minioutil "crf-mold/common/minio"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	MOLD                  = "mold"
	SPARE                 = "spare"
	SPARE_REQUEST         = "spare_request"
	MOLD_BOM              = "mold_bom"
	PRODUCT_RESUME        = "product_resume"
	SPARE_RESUME          = "spare_resume"
	STEREOSCOPIC_LOCATION = "stereoscopic_location"
)

type MinIOFileVO struct {
	BucketName string `json:"bucketName"`
	FileKey    string `json:"fileKey"`
	FileName   string `json:"fileName"`
}

var ExcelTemplateEnum = map[string]*MinIOFileVO{
	MOLD:          {minioutil.BucketName, "模具台账.xlsx", "模具台账.xlsx"},
	SPARE:         {minioutil.BucketName, "备件基础资料.xlsx", "备件基础资料.xlsx"},
	SPARE_REQUEST: {minioutil.BucketName, "备件库存.xlsx", "备件库存.xlsx"},
	MOLD_BOM:      {minioutil.BucketName, "模具BOM.xlsx", "模具BOM.xlsx"},
}

/**
 * 生成一个excel文件名称
 */
func GenerateExcelFile(moduler string) (fileName string) {
	return moduler + "_" + base.UUID() + ".xlsx"
}

// @Tags Excel模板
// @Summary 下载导入模板
// @Accept json
// @Produce json
// @Param module query string true "模块：模具mold，备件spare，库存spare_request,模具BOM:mold_bom"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /template [get]
func DownloadTemplate(c *gin.Context) {
	module := c.Query("module")
	if module == "" {
		panic(base.ParamsErrorN())
	}

	res := ExcelTemplateEnum[module]
	if res == nil {
		panic(base.ParamsErrorN())
	}

	c.JSON(http.StatusOK, base.Success(res))
}
