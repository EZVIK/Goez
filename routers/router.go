package routers

import (
	"Goez/handler"
	"Goez/middleware"
	"Goez/pkg/upload"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() *gin.Engine {

	r := gin.Default()

	r.Use(middleware.CORS)

	// 图片访问
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))

	// -------- 对外接口
	api := r.Group("/api/v1")

	// 用户
	api.POST("/login", handler.Login)       // 用户登陆
	api.POST("/register", handler.Register) // 用户注册

	// 文章
	api.GET("/article/:id", handler.QueryArticleById)                // 根据id获取文章详情
	api.GET("/articles/:field/:value", handler.QueryArticleByFields) // 带参数模糊查询文章列表
	api.GET("/articles", handler.QueryArticles)                      // 带参数模糊查询文章列表

	// JWT校验
	api.Use(middleware.JWT)

	// 标签
	api.GET("/tag/:id", handler.GetTagById)
	api.POST("/tag", handler.AddTag)
	api.PATCH("/tag", handler.UpdateTagById)
	api.DELETE("/tag", handler.DeleteTag)
	api.GET("/relatedTags", handler.Recommand)
	api.GET("/articleTags/:id", handler.GetArticleTag)

	// -------- 后台管理接口
	auth := r.Group("/auth/api/v1")

	//
	auth.POST("/article", handler.AddArticle)
	auth.PATCH("/article", handler.UpdateArticle)
	auth.DELETE("/article", handler.DeleteArticle)
	auth.GET("/article/:id", handler.QueryArticleByIdAuth)
	auth.POST("/upload", handler.UploadImage)
	auth.POST("/article/export", handler.ExportArticleData)

	return r
}
