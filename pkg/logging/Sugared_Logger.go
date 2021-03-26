package logging

import (
	"Goez/pkg/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
)

var sugarLogger *zap.SugaredLogger

//func main() {
//	InitLogger()
//	defer sugarLogger.Sync()
//	simpleHttpGet("www.google.com")
//	simpleHttpGet("http://159.75.82.148:8000/api/v1/articles?pageSize=50")
//}

func Setup() {
	writeSyncer := getLogWriter()
	encoder := getEncode()

	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())

	sugarLogger = logger.Sugar()
}

func GetLogger() *zap.SugaredLogger {
	return sugarLogger
}

func getEncode() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 设置ISO08601 时间编码
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	file := config.AppSetting.RuntimeRootPath + config.AppSetting.LogSavePath + config.AppSetting.LogSaveName
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
		LocalTime:  true,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func Info(info string, args interface{}) {
	sugarLogger.Infof(info, args)
}

func Error(err string, args interface{}) {
	sugarLogger.Errorf(err, args)
}

func simpleHttpGet(url string) {
	sugarLogger.Debugf("Trying to hit GET request for %s", url)
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Errorf("Error fetching URL %s : Error = %s", url, err)
	} else {
		sugarLogger.Infof("Success! statusCode = %s for URL %s", resp.Status, url)
		resp.Body.Close()
	}
}
