package middleware

import (
	"Goez/pkg/e"
	"Goez/pkg/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JWT(c *gin.Context) {
	var code int
	var data interface{}

	code = e.SUCCESS
	token := c.GetHeader("token")

	if token == "" {
		code = e.INVALID_PARAMS
	} else {
		_, err := utils.ParseToken(token)
		if err != nil {
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			default:
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			}
		}
	}

	if code != e.SUCCESS {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": data,
		})

		c.Abort()
		return
	}

	c.Next()
}
