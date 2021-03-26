package middleware

import (
	"Goez/pkg/logging"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"strings"
	"time"
)

// LogConfig 日志配置
type LogConfig struct {
	SkipPaths map[string]bool
}

type ServiceLog struct {
	Reqid  string
	Ip     string
	Custom string
	Msg    string
}

func getBasicHTTPReqInfo(c *gin.Context) map[string]string {
	ip := c.ClientIP()

	return map[string]string{
		"uri":         c.Request.RequestURI,
		"http-method": c.Request.Method,
		"req-id":      c.Request.Header.Get("x-request-id"),
		"req-ip":      ip,
	}
}

func getBody2JSONStr(body io.ReadCloser) (string, io.ReadCloser, error) {

	//buf := new(bytes.Buffer)
	//buf.ReadFrom(body)
	//buf.Reset()
	//
	//s := buf.String()
	//return s

	b, err := ioutil.ReadAll(body)

	if err == nil {
		stringReader := bytes.NewReader(b)

		stringReadCloser := ioutil.NopCloser(stringReader)

		return string(b), stringReadCloser, nil
	}

	return "", nil, err
}

// 请求LOG 中间件
func LoggingMiddleware(conf LogConfig) func(*gin.Context) {

	return func(c *gin.Context) {

		// 不需要记录的 URL
		if conf.SkipPaths[c.Request.URL.RawPath] {
			return
		}

		// 获取请求信息
		infoMap := getBasicHTTPReqInfo(c)

		// 开始时间
		startTime := time.Now()

		reqBody := ""
		// 获取请求body
		if body, rc, err := getBody2JSONStr(c.Request.Body); err != nil {
			logging.GetLogger().Errorf("getBody2JSONStr error: %v", err)
		} else {
			c.Request.Body = rc
			reqBody = strings.Replace(body, "\n", "", -1)
		}

		svcLog := ServiceLog{
			Reqid: infoMap["req-id"],
			Ip:    infoMap["req-ip"],
			Custom: fmt.Sprintf(
				"uri: %v %v, req-body: %v",
				infoMap["http-method"],
				infoMap["uri"],
				reqBody,
			),
			Msg: fmt.Sprintf("start calling at: %v", startTime.Format("2006-01-02T15:04:05-0700")),
		}

		logging.GetLogger().Info(svcLog)

		c.Next() // 处理请求

		latencyTime := time.Now().Sub(startTime) // 执行时间
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		statusCode := c.Writer.Status()

		svcLog.Custom = fmt.Sprintf(
			"uri: %v %v, http-status: %v",
			infoMap["http-method"],
			infoMap["uri"],
			statusCode,
		)
		svcLog.Msg = fmt.Sprintf("called, consuming time: %v", latencyTime)
		logging.GetLogger().Info(svcLog)
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func ginBodyLogMiddleware(c *gin.Context) {
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw
	c.Next()
	statusCode := c.Writer.Status()
	if statusCode >= 400 {
		//ok this is an request with error, let's make a record for it
		// now print body (or log in your preferred way)
		fmt.Println("Response body: " + blw.body.String())
	}
}
