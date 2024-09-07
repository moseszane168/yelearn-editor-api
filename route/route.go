package route

import (
	"crf-mold/controller/editor/coursepacks"
	"crf-mold/controller/editor/courses"
	"crf-mold/controller/editor/statements"
	"crf-mold/controller/file"
	"github.com/gin-gonic/gin"
)

func InitRoute(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		// 课程包
		v1.POST("/editor/course-pack", coursepacks.CreateCoursePack)   // 新增课程包
		v1.DELETE("/editor/course-pack", coursepacks.DeleteCoursePack) // 删除课程包
		v1.PUT("/editor/course-pack", coursepacks.UpdateCoursePack)    // 更新课程包
		v1.GET("/editor/course-packs", coursepacks.ListCoursePack)     // 课程包列表
		v1.GET("/editor/course-pack", coursepacks.OneCoursePack)       // 课程包查看
		v1.POST("/file", file.UploadFile)                              // 上传课程包封面

		// 课程
		v1.POST("/editor/course", courses.CreateCourse)   // 新增课程
		v1.DELETE("/editor/course", courses.DeleteCourse) // 删除课程
		v1.PUT("/editor/course", courses.UpdateCourse)    // 更新课程
		v1.GET("/editor/courses", courses.ListCourse)     // 课程列表
		v1.GET("/editor/course", courses.OneCourse)       //课程查看

		// 句子
		v1.POST("/editor/statement", statements.CreateStatement)      // 新增句子
		v1.DELETE("/editor/statement", statements.DeleteStatement)    // 删除句子
		v1.PUT("/editor/statement", statements.UpdateStatement)       // 更新句子
		v1.GET("/editor/statements", statements.ListStatement)        // 句子列表
		v1.GET("/editor/statement", statements.OneStatement)          // 句子查看
		v1.POST("/editor/split-statement", statements.SplitStatement) // 拆分句子

	}
}
