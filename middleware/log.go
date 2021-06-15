package middleware
//
//import (
//	"fmt"
//	"github.com/gin-gonic/gin"
//)
//
//// CostumerLog 自定义日志输出格式
//func CostumerLog() gin.HandlerFunc {
//	// LoggerWithFormatter 中间件会将日志写入 gin.DefaultWriter
//	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
//		return fmt.Sprintf("%s |\t%s |\t%d |\t%s |\t\"%s\" |\t%s |\t%s |\t%s |\t%s \n",
//			param.TimeStamp.Format("2006-01-02 03:04:05"),
//			param.ClientIP,
//			param.StatusCode,
//			param.Method,
//			param.Path,
//			param.Request.Proto,
//			param.Latency,
//			param.Request.UserAgent(),
//			param.ErrorMessage,
//		)
//	})
//}
