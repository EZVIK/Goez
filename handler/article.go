package handler

import (
	"Goez/handler/dto"
	"Goez/models"
	"Goez/pkg/app"
	"Goez/pkg/e"
	"Goez/pkg/utils"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// article

// ----------- 前台接口 ----------- //

// 根据ID 查询文章详情 /article/:id
func QueryArticleById(c *gin.Context) {
	appG := app.Gin{C: c}
	code := e.SUCCESS
	aid := c.Param("id")
	i, err := strconv.Atoi(aid)

	if err != nil {
		code = e.ERROR
	}

	art := models.Article{ID: i}

	// 查询文章
	a, err := art.GetArticlesById()

	if err != nil {
		code = e.INVALID_PARAMS
		status := http.StatusBadRequest
		if err.Error() == "record not found" {
			status = http.StatusNotFound
		}
		appG.Response(status, e.NOT_FOUND, err)
		return
	}

	// 获取用户id
	token := c.GetHeader("token")
	claims, err := utils.ParseToken(token)
	if err != nil && claims != nil {
		appG.Response(http.StatusBadRequest, code, err)
		return
	}

	// 添加浏览记录
	r := models.Record{UserId: claims.UserId, ArticleId: i}

	if _, err := models.AddRecord(r); err != nil {
		appG.Response(http.StatusBadRequest, code, err)
		return
	}

	c.JSON(code, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": a,
	})
}

//	带参数模糊 查询文章列表 /articles/:field/:value
func QueryArticleByFields(c *gin.Context) {
	appG := app.Gin{C: c}

	pageIndex, pageSize := models.PageChecker(c.Query("page"), "20")

	ap := dto.ArticleSearchParams{
		c.Param("field"),
		c.Param("value"),
		pageIndex, pageSize,
	}

	if errParams := validator.New().Struct(ap); errParams != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	art := models.Article{}

	al, count, err := art.GetArticleByFields(ap.Field, ap.Value, 0, 0)

	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	resp := map[string]interface{}{
		"list":  al,
		"count": count,
		"page":  pageIndex,
	}

	appG.Response(http.StatusOK, e.SUCCESS, resp)
}

func QueryArticles(c *gin.Context) {
	appG := app.Gin{C: c}

	art := models.Article{}

	pageIndex, pageSize := models.PageChecker(c.Query("page"), c.Query("pageSize"))

	al, count, err := art.GetArticleByFields("", "", pageIndex, pageSize)

	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	resp := map[string]interface{}{
		"list":  al,
		"count": count,
	}

	appG.Response(http.StatusOK, e.SUCCESS, resp)
}

// 推荐文章列表
func Recommand(c *gin.Context) {
	appG := app.Gin{C: c}

	// 获取用户id
	token := c.GetHeader("token")
	claims, err := utils.ParseToken(token)

	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	// 根据用户ID 查询浏览过的文章 相关的标签
	userId := claims.UserId
	tags, err := models.GetTagsFromRecord(userId)

	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	var buffer bytes.Buffer
	// 合并标签 为字符串
	for i, k := range tags {
		buffer.WriteString(k.Name)
		if i != len(tags)-1 {
			buffer.WriteString(" ")
		}
	}

	reqData := map[string]string{
		"id":   strconv.Itoa(userId),
		"tags": buffer.String(),
	}

	// 请求匹配推荐文章
	respData, err := resty.New().R().
		SetHeader("Content-Type", "multipart/form-data").
		SetFormData(reqData).
		Post("http://localhost:5000/recomman")

	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	reqDa := TagText{}
	// 推荐文章列表
	if err := json.Unmarshal(respData.Body(), &reqDa); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, reqDa.Data)
}

type TagText struct {
	Data []string `json:"data"`
}

// ----------- 后台管理 ----------- //

//  查询文章详情 /article/:id TODO 直接查询数据库 不使用缓存
func QueryArticleByIdAuth(c *gin.Context) {
	appG := app.Gin{C: c}

	aid := c.Param("id")
	i, err := strconv.Atoi(aid)

	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	art := models.Article{ID: i}

	a, err := art.GetArticlesById()

	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, a)
}

//	查询文章列表
func QueryArticleList(c *gin.Context) {
	appG := app.Gin{C: c}

	// TODO QUERY PARAMS GetArticleTitleList(PARAMS)

	al, err := models.Article{}.GetArticleTitleList()

	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, al)
}

// 添加文章
func AddArticle(c *gin.Context) {
	appG := app.Gin{C: c}

	code := e.SUCCESS

	aa := dto.AddArticleParams{}

	if err := appG.C.ShouldBind(&aa); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, err)
		return
	}

	// TODO 处理标题规范 不大于30字符
	// TODO desc 不大于50字符
	// TODO tags长度大于 1

	// 获取用户id
	token := c.GetHeader("token")
	claims, err := utils.ParseToken(token)

	if err != nil {
		code = e.ERROR
	}

	var ts []models.Tag
	art := models.Article{0, aa.Title, aa.Desc, aa.Content, 0, claims.UserId, ts, gorm.Model{}}

	// loop tags check if insert
	for _, temp_tag := range aa.Tags {

		if !checkTagName(temp_tag) {
			continue
		}

		checkTag := models.Tag{Name: temp_tag}

		t, err := checkTag.GetTag()
		if err != nil {
			appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, err)
			return
		}
		// 添加到文章标签列表
		ts = append(ts, t)
	}

	art.Tags = ts

	a, err := art.AddArticle()

	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR_ADD_ARTICLE_FAIL, "AddArticle fails :"+art.Title+"be posted.1"+err.Error())
		return
	}

	c.JSON(code, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": a.ID,
	})
}

// 更新文章
func UpdateArticle(c *gin.Context) {

}

// 删除文章（软删除）
func DeleteArticle(c *gin.Context) {
	appG := app.Gin{C: c}

	art := models.Article{}

	if err := c.ShouldBind(&art); err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	if err := art.DeleteArticle(); err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR_DELETE_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, "ok")
}

func checkTagName(temp_tag string) bool {

	if temp_tag == "" {
		return false
	}

	return true
}
