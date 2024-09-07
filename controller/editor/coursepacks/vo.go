package coursepacks

import (
	"crf-mold/base"
	"time"
)

type CreateCoursePackVO struct {
	IsFree      bool      `json:"isFree"`                         //
	Order       int       `json:"order"`                          //
	CreatedAt   time.Time `json:"createdAt"`                      //
	UpdatedAt   time.Time `json:"updatedAt"`                      //
	CreatorId   string    `json:"creatorId"`                      //
	ShareLevel  string    `json:"shareLevel"`                     //
	Cover       string    `json:"cover" binding:"required"`       // 封面
	Title       string    `json:"title" binding:"required"`       // 标题
	Description string    `json:"description" binding:"required"` // 描述
	Id          string    `json:"id" binding:"required"`          // 主键id
} // @name CreateCoursePackVO

type UpdateCoursePackVO struct {
	IsFree      bool      `json:"isFree"`                         //
	Order       int       `json:"order"`                          //
	CreatedAt   time.Time `json:"createdAt"`                      //
	UpdatedAt   time.Time `json:"updatedAt"`                      //
	CreatorId   string    `json:"creatorId"`                      //
	ShareLevel  string    `json:"shareLevel"`                     //
	Cover       string    `json:"cover" binding:"required"`       // 封面
	Title       string    `json:"title" binding:"required"`       // 标题
	Description string    `json:"description" binding:"required"` // 描述
	Id          string    `json:"id" binding:"required"`          // 主键id
} // @name UpdateCoursePackVO

type PageCoursePackVO struct {
	base.PageVO
	IsFree      bool      `json:"isFree"`                         //
	Order       int       `json:"order"`                          //
	CreatedAt   time.Time `json:"createdAt"`                      //
	UpdatedAt   time.Time `json:"updatedAt"`                      //
	CreatorId   string    `json:"creatorId"`                      //
	ShareLevel  string    `json:"shareLevel"`                     //
	Cover       string    `json:"cover" binding:"required"`       // 封面
	Title       string    `json:"title" binding:"required"`       // 标题
	Description string    `json:"description" binding:"required"` // 描述
	Id          string    `json:"id" binding:"required"`          // 主键id
} // @name PageCoursePackVO

type PageCoursePackOutputVO struct {
	IsFree      bool      `json:"isFree"`                         //
	Order       int       `json:"order"`                          //
	CreatedAt   time.Time `json:"createdAt"`                      //
	UpdatedAt   time.Time `json:"updatedAt"`                      //
	CreatorId   string    `json:"creatorId"`                      //
	ShareLevel  string    `json:"shareLevel"`                     //
	Cover       string    `json:"cover" binding:"required"`       // 封面
	Title       string    `json:"title" binding:"required"`       // 标题
	Description string    `json:"description" binding:"required"` // 描述
	Id          string    `json:"id" binding:"required"`          // 主键id
} // @name PageCoursePackOutputVO

type CoursePackOutVO struct {
	IsFree      bool      `json:"isFree"`                         //
	Order       int       `json:"order"`                          //
	CreatedAt   time.Time `json:"createdAt"`                      //
	UpdatedAt   time.Time `json:"updatedAt"`                      //
	CreatorId   string    `json:"creatorId"`                      //
	ShareLevel  string    `json:"shareLevel"`                     //
	Cover       string    `json:"cover" binding:"required"`       // 封面
	Title       string    `json:"title" binding:"required"`       // 标题
	Description string    `json:"description" binding:"required"` // 描述
	Id          string    `json:"id" binding:"required"`          // 主键id
} // @name CoursePackOutVO
