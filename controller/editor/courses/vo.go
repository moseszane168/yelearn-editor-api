package courses

import (
	"crf-mold/base"
	"time"
)

type CreateCourseVO struct {
	UpdatedAt    time.Time `json:"updatedAt"`                       //
	Order        int       `json:"order"`                           //
	CreatedAt    time.Time `json:"createdAt"`                       //
	Video        string    `json:"video"`                           //
	Id           string    `json:"Id"`                              // 主键id
	CoursePackId string    `json:"coursePackId" binding:"required"` // 课程id
	Title        string    `json:"Title" binding:"required"`        // 封面
	Description  string    `json:"description"`                     // 描述
} // @name CreateCourseVO

type UpdateCourseVO struct {
	UpdatedAt    time.Time `json:"updatedAt"`             //
	Order        int       `json:"order"`                 //
	CreatedAt    time.Time `json:"createdAt"`             //
	Video        string    `json:"video"`                 //
	Id           string    `json:"Id" binding:"required"` // 主键id
	CoursePackId string    `json:"CoursePackId"`          // 课程id
	Title        string    `json:"Title"`                 // 封面
	Description  string    `json:"description"`           // 描述
} // @name UpdateCourseVO

type PageCourseVO struct {
	base.PageVO
	UpdatedAt    time.Time `json:"updatedAt"`                       //
	Order        int       `json:"order"`                           //
	CreatedAt    time.Time `json:"createdAt"`                       //
	Video        string    `json:"video"`                           //
	Id           string    `json:"Id" binding:"required"`           // 主键id
	CoursePackId string    `json:"CoursePackId" binding:"required"` // 课程id
	Title        string    `json:"Title" binding:"required"`        // 封面
	Description  string    `json:"description" binding:"required"`  // 描述
} // @name PageCourseVO

type PageCourseOutputVO struct {
	UpdatedAt    time.Time `json:"updatedAt"`                       //
	Order        int       `json:"order"`                           //
	CreatedAt    time.Time `json:"createdAt"`                       //
	Video        string    `json:"video"`                           //
	Id           string    `json:"Id" binding:"required"`           // 主键id
	CoursePackId string    `json:"CoursePackId" binding:"required"` // 课程id
	Title        string    `json:"Title" binding:"required"`        // 封面
	Description  string    `json:"description" binding:"required"`  // 描述
} // @name PageCourseOutputVO

type CourseOutVO struct {
	UpdatedAt    time.Time `json:"updatedAt"`                       //
	Order        int       `json:"order"`                           //
	CreatedAt    time.Time `json:"createdAt"`                       //
	Video        string    `json:"video"`                           //
	Id           string    `json:"Id" binding:"required"`           // 主键id
	CoursePackId string    `json:"CoursePackId" binding:"required"` // 课程id
	Title        string    `json:"Title" binding:"required"`        // 封面
	Description  string    `json:"description" binding:"required"`  // 描述
} // @name CourseOutVO
