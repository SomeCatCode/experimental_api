package application

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func LogFormatDetail(param gin.LogFormatterParams) string {
	return fmt.Sprintf("[GIN] %s | %3d | %13v | %15s |%-7s %#v | UA: %s | Header: %v\n",
		param.TimeStamp.Format(time.RFC3339),
		param.StatusCode,
		param.Latency,
		param.ClientIP,
		param.Method,
		param.Path,
		param.Request.UserAgent(),
		param.Request.Header, // oder param.Request.Header
	)
}

func LogFormatShort(param gin.LogFormatterParams) string {
	return fmt.Sprintf("[GIN] %s | %3d | %13v | %15s |%-7s %#v\n",
		param.TimeStamp.Format(time.RFC3339),
		param.StatusCode,
		param.Latency,
		param.ClientIP,
		param.Method,
		param.Path,
	)
}
