package main

import (
	"Goez/models"
	"Goez/pkg/config"
	"Goez/pkg/logging"
	"Goez/routers"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func init() {
	config.Setup()
	models.Setup()
	logging.Setup()
}

func main() {
	gin.SetMode(config.ServerSetting.RunMode)

	routersInit := routers.InitRouter()
	readTimeout := config.ServerSetting.ReadTimeout
	writeTimeout := config.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", config.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	logging.Info("start app server listening %s", endPoint)

	server.ListenAndServe()
}
