/**
 * 分页工具类，统一分页处理
 */

package base

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const DEFAULT_CURRENT_PAGE int64 = 1
const DEFAULT_SIZE int64 = 10

type PageVO struct {
	CurrentPage int64 `json:"currentPage" form:"currentPage"` // 当前页
	Size        int64 `json:"size" form:"size"`               // 每页条数
}

func (p *PageVO) GetCurrentPage() int64 {
	if p.CurrentPage == 0 {
		return DEFAULT_CURRENT_PAGE
	}
	return p.CurrentPage
}

func (p *PageVO) GetSize() int64 {
	if p.CurrentPage == 0 {
		return DEFAULT_SIZE
	}
	return p.Size
}

type BasePage struct {
	Size        int64       `json:"size"`
	CurrentPage int64       `json:"currentPage"`
	TotalPage   int64       `json:"totalPage"`
	TotalRows   int64       `json:"totalRows"`
	List        interface{} `json:"list"`
}

type CountVO struct {
	Count int64
}

/**
 * 创建一个统一的分页对象返回
 */
func NewBasePage(current, size, totalRows int64, list interface{}) *BasePage {
	var page BasePage

	pageNum := totalRows / size
	if totalRows%size == 0 {
		page.TotalPage = pageNum
	} else {
		page.TotalPage = pageNum + 1
	}

	page.CurrentPage = current
	page.Size = size
	page.TotalRows = totalRows
	page.List = list
	return &page
}

/**
 * 拓展gorm，实现分页功能,不支持Having中带有别名的条件
 */
func Page(tx *gorm.DB, result interface{}, current, size int64) *BasePage {
	var count int64
	tx.Count(&count)
	tx.Limit(int(size)).Offset(int((current - 1) * size)).Find(result)
	return NewBasePage(current, size, count, result)
}

/**
 * 支持Having中带有别名的条件
 */
func PageLowLevelCountSQL(tx *gorm.DB, result interface{}, current, size int64) *BasePage {
	oriSql := tx.ToSQL(func(tx *gorm.DB) *gorm.DB { return tx.Find(&[]interface{}{}) })

	countSql := "select count(1) as count from (" + oriSql + ") t"
	var countVo CountVO
	tx.Raw(countSql).Scan(&countVo)

	tx.Limit(int(size)).Offset(int((current - 1) * size)).Find(result)
	return NewBasePage(current, size, countVo.Count, result)
}

/**
 * 原始SQL分页
 */
func PageWithRawSQL(tx *gorm.DB, result interface{}, current, size int64, rawSql string, values ...interface{}) *BasePage {
	countSql := "select count(1) as count from (" + rawSql + ") t"
	var countVo CountVO
	tx.Raw(countSql, values...).Scan(&countVo)

	pageSql := rawSql + fmt.Sprintf(" limit %d,%d", int((current-1)*size), size)
	tx.Raw(pageSql, values...).Scan(result)

	return NewBasePage(current, size, countVo.Count, result)
}

/**
 * 从URL的查询参数中取分页对象
 */
func GetPageParams(c *gin.Context) (int64, int64) {
	currentPageParam := c.Query("currentPage")
	sizeParam := c.Query("size")

	var current, size int64
	if currentPageParam == "" {
		current = DEFAULT_CURRENT_PAGE
	} else {
		c, err := strconv.Atoi(currentPageParam)
		if err != nil {
			panic(ParamsErrorN())
		}
		current = int64(c)
	}
	if sizeParam == "" {
		size = DEFAULT_SIZE
	} else {
		s, err := strconv.Atoi(sizeParam)
		if err != nil {
			panic(ParamsErrorN())
		}
		size = int64(s)
	}

	return current, size
}
