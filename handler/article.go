package handler

import (
	"Goez/handler/dto"
	"Goez/models"
	"Goez/pkg/app"
	"Goez/pkg/e"
	"Goez/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
)

// article

// ----------- 前台接口 ----------- //

// 根据ID 查询文章详情 /article/:id
func QueryArticleById(c *gin.Context)  {
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

	// 添加浏览记录
	r := models.Record{UserId: claims.UserId, ArticleId: i}

	if _, err := models.AddRecord(r); err != nil {
		appG.Response(http.StatusBadRequest, code, err)
		return
	}

	c.JSON(code, gin.H{
		"code": code,
		"msg" : e.GetMsg(code),
		"data": a,
	})
}

//	带参数模糊 查询文章列表 /articles/:field/:value
func QueryArticleByFields(c *gin.Context)  {
	appG := app.Gin{C: c}

	ap := dto.ArticleSearchParams{
		c.Param("field"),
		c.Param("value"),
	}

	if errParams := validator.New().Struct(ap); errParams != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	art := models.Article{}

	al, count, err := art.GetArticleByFields(ap.Field, ap.Value)

	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	resp := map[string]interface{}{
		"list": al,
		"count": count,
	}

	appG.Response(http.StatusOK, e.SUCCESS, resp)
}


// ----------- 后台管理 ----------- //

//  查询文章详情 /article/:id TODO 直接查询数据库 不使用缓存
func QueryArticleByIdAuth(c *gin.Context)  {
	aid := c.Param("id")
	i, err := strconv.Atoi(aid)
	code := e.SUCCESS

	if err != nil {
		code = e.ERROR
	}

	art := models.Article{ID: i}

	a, err := art.GetArticlesById()

	if err != nil {
		code = e.INVALID_PARAMS
	}

	c.JSON(code, gin.H{
		"code": code,
		"msg" : e.GetMsg(code),
		"data": a,
	})
}

//	查询文章列表
func QueryArticleList(c *gin.Context)  {
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
func AddArticle(c *gin.Context)  {
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

	art := models.Article{0, aa.Title, aa.Desc, aa.Content, 0,claims.UserId, &ts, models.NewModelTime()}

	if _, alen, err := art.GetArticleByFields("title", art.Title); err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR_GET_ARTICLES_FAIL, "GetTagsById error:" + err.Error())
		return
	} else if alen > 0 {
		appG.Response(http.StatusBadRequest, e.ERROR_GET_ARTICLES_FAIL, "AddArticle fails :"+ art.Title +"be posted.")
		return
	}

	// loop tags check if insert
	for _, temp_tag := range aa.Tags {

		if !checkTagName(temp_tag) {
			continue
		}

		checkTag := models.Tag{Name: temp_tag}

		t, err := checkTag.GetTag()

		// 查询成功
		if err != nil {

			// 标签未添加
			if err.Error() == "record not found" {
				if newTag, err_add := checkTag.AddTag(); err_add == nil {
					// 添加成功
					t = newTag
				} else {
					// 添加失败
					appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "GetTagsById error:" + err_add.Error())
					return
				}
			} else {
				appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "GetTagsById error:" + err.Error())
				return
			}
		}

		// 添加到文章标签列表
		ts = append(ts, t)
	}

	art.Tags = &ts

	a, err := art.AddArticle()

	c.JSON(code, gin.H{
		"code": code,
		"msg" : e.GetMsg(code),
		"data": a.ID,
	})
}

// 更新文章
func UpdateArticle(c *gin.Context)  {

}

// 删除文章（软删除）
func DeleteArticle(c *gin.Context)  {
	code := e.SUCCESS
	id := c.PostForm("id")
	i,err := strconv.Atoi(id)
	art := models.Article{ID: i}

	if err != nil {
		code = e.INVALID_PARAMS
	} else {
		ok, err := art.DeleteArticle();
		if err != nil && ok{
			code = e.ERROR_DELETE_ARTICLE_FAIL
		}
	}

	c.JSON(code, gin.H{
		"code": code,
		"msg" : e.GetMsg(code),
		"data"  : "ok",
	})
}

func checkTagName(temp_tag string) (bool) {

	if temp_tag == "" {
		return false
	}

	return true
}