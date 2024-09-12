package middleware

import (
	"github.com/fanxiangqing/web-base/lib/utils"
	"github.com/fanxiangqing/web-base/lib/utils/types"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var NotLogPath = types.NewStrSet(false, "/favicon.ico")    //NotLogPath  不用记录日志的路径
var NotReqLogPath = types.NewStrSet(false, "/favicon.ico") //NotReqDataPath  日志不用记录请求参数的路径

// LoggingJson 带field的json记录请求日志
func LoggingJson() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		if !ShouldLog(ctx) {
			return
		}

		// 开始时间
		start := time.Now()
		// 处理请求
		ctx.Next()
		logId := utils.GetLogId(ctx)
		// 结束时间
		end := time.Now()
		//执行时间
		latency := end.Sub(start)
		clientIP := ctx.ClientIP()
		method := ctx.Request.Method
		statusCode := ctx.Writer.Status()
		code := ctx.GetInt("code")  // 2023-10-27, 记录接口业务状态
		msg := ctx.GetString("msg") // 2023-10-27, 记录接口业务提示信息
		rawData := ctx.Request.URL.RawQuery
		form := ctx.Request.PostForm.Encode()

		if !ShouldLogReq(ctx) {
			//无意义的加密数据等路径，对于查日志无用
			//可能关心接口的响应状态值
			rawData = ""
			form = ""

		}
		// Logrus 鼓励用户通过日志字段记录结构化日志
		// 结构化日志有利于工具提取并分析日志。
		logrus.WithFields(logrus.Fields{
			"logType":  "access",
			"clientIp": clientIP,
			"status":   statusCode,
			"duration": latency.Seconds(), //统一单位秒
			"request":  ctx.Request.URL.Path,
			"query":    rawData,
			"method":   method,
			"form":     form,
			"apiCode":  code,
			"apiMsg":   msg,
			"logId":    logId,
		}).Infof("access")
	}
}

func ShouldLog(ctx *gin.Context) bool {
	return !NotLogPath.Contains(ctx.Request.URL.Path)
}

func ShouldLogReq(ctx *gin.Context) bool {
	return !NotReqLogPath.Contains(ctx.Request.URL.Path)
}
