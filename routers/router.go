package routers

import (
	"Goez/handler"
	"Goez/middleware/jwt"
	"Goez/pkg/upload"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() *gin.Engine {

	r := gin.Default()

	// 图片访问
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))

	// -------- 对外接口
	api := r.Group("/api/v1")

	// 用户
	api.POST("/login", handler.Login)											// 用户登陆
	api.POST("/register", handler.Register)										// 用户注册

	// 文章
	api.GET("/article/:id", handler.QueryArticleById)							// 根据id获取文章详情
	api.GET("/articles/:field/:value", handler.QueryArticleByFields)				// 带参数模糊查询文章列表
	//api.GET("/articles", handler.QueryArticleList)
	// api.POST("comment", handler.AddComment)

	// 标签
	api.GET("/tag/:id", jwt.JWT(),handler.GetTagById)
	api.POST("/tag", jwt.JWT(),   handler.AddTag)
	api.PATCH("/tag", jwt.JWT(),  handler.UpdateTagById)
	api.DELETE("/tag", jwt.JWT(), handler.DeleteTag)

	// -------- 后台管理接口
	auth := r.Group("/auth/api/v1")
	auth.POST("/article", jwt.JWT(), handler.AddArticle)
	auth.PATCH("/article", jwt.JWT(), handler.UpdateArticle)
	auth.DELETE("/article", jwt.JWT(), handler.DeleteArticle)
	auth.GET("/article/:id", jwt.JWT(), handler.QueryArticleByIdAuth)
	auth.POST("/upload", jwt.JWT(), handler.UploadImage)

	return r
}