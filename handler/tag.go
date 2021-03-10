package handler

import (
	"Goez/handler/dto"
	"Goez/models"
	"Goez/pkg/app"
	"Goez/pkg/e"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddTag(c *gin.Context)  {
	appG := app.Gin{C: c}

	atp := new(dto.AddTagsParams)

	if err := appG.C.ShouldBind(&atp); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, err)
		return
	}

	err := errors.New("")

	if atp.Name == "" {
		var tags []models.Tag
		for _, k := range atp.Names {
			tags = append(tags, models.Tag{Name: k})
		}
		err = models.AddTags(tags)
	} else {
		_, err = models.Tag{Name: atp.Name}.AddTag()
	}

	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, err)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func GetTagById(c *gin.Context)  {
	appG := app.Gin{C: c}

	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

func UpdateTagById(c *gin.Context)  {
	appG := app.Gin{C: c}

	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

func DeleteTag(c *gin.Context)  {
	appG := app.Gin{C: c}

	appG.Response(http.StatusOK, e.SUCCESS, nil)

}