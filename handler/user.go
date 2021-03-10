package handler

import (
	"Goez/handler/dto"
	"Goez/models"
	"Goez/pkg/app"
	"Goez/pkg/e"
	"Goez/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 登陆
func Login(c *gin.Context)  {
	appG := app.Gin{C: c}
	code := e.SUCCESS

	loginParams := dto.LoginParams{}
	if err := appG.C.ShouldBind(&loginParams); err != nil {
		code := e.INVALID_PARAMS
		appG.Response(http.StatusOK, code, nil)
		return
	}

	username := loginParams.Username
	password := loginParams.Password

	u, err := models.CheckAuth(username, password)

	if err != nil || u.ID == 0{
		code := e.INVALID_PARAMS
		appG.Response(http.StatusOK, code, nil)
		return
	}

	token, err := utils.GenerateToken(username, password, u.ID)

	if err != nil {
		code = e.ERROR_AUTH
	}

	appG.Response(http.StatusOK, code, token)
}

// 注册
func Register(c *gin.Context)  {
	reg := dto.RegisterParams{}
	appG := app.Gin{C: c}

	if err := c.ShouldBind(&reg); err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	if reg.Password != reg.RePassword {
		appG.Response(http.StatusOK, e.ERROR_AUTH_REGISTER_PASSWORD, nil)
		return
	}

	id, err := models.Register(reg.Username, reg.Password, reg.Nickname)

	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR_AUTH_REGISTER_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, id)
}

func GetUserInfo(c *gin.Context) {

}
