package logging

import (
	"Goez/pkg/config"
	"fmt"
	"time"
)

func getLogFilePath() string {
	return fmt.Sprintf("%s%s", config.AppSetting.RuntimeRootPath, config.AppSetting.LogSavePath)
}

func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s",
		config.AppSetting.LogSaveName,
		time.Now().Format(config.AppSetting.TimeFormat),
		config.AppSetting.LogFileExt,
	)
}
