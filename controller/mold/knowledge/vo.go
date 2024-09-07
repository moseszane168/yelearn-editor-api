package mold

import "crf-mold/base"

//
// 模具知识库
//

type CreateMoldKnowledgeVO struct {
	Name    string        `json:"name" binding:"required"`    // 知识库名称
	Content string        `json:"content" binding:"required"` // 知识库内容
	Group   string        `json:"group" binding:"required"`   // 知识库组
	Style   string        `json:"style" binding:"required"`   // 书皮样式
	Files   base.FileList `json:"files"`                      // 上传保存文件
} // @name CreateMoldKnowledgeVO

type UpdateMoldKnowledgeVO struct {
	ID      int64         `json:"id" binding:"required"`    // 知识库ID
	Name    string        `json:"name" binding:"required"`  // 知识库名称
	Content string        `json:"content"`                  // 知识库内容
	Group   string        `json:"group" binding:"required"` // 知识库组
	Style   string        `json:"style"`                    // 书皮样式
	Files   base.FileList `json:"files"`                    // 上传保存文件
} // @name UpdateMoldKnowledgeVO

type KnowledgeSearchVO struct {
	ID         int64  `json:"id"`         // 知识库ID
	Name       string `json:"name"`       // 知识库名称
	Content    string `json:"content"`    // 知识库内容
	Group      string `json:"group"`      // 知识库组
	CreatedBy  string `json:"createdBy"`  // 创建人
	GmtUpdated string `json:"gmtUpdated"` // 更新时间
} // @name KnowledgeSearchVO
