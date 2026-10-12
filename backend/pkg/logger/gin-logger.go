package logger

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	// LoggerConfig is the config for the logger
	LoggerConfig = gin.LoggerConfig{
		Formatter: DefaultLogFormatter,
		Skip: func(c *gin.Context) bool {
			return strings.HasPrefix(c.Request.URL.Path, "/internal")
		},
	}
)

// ContextKeyUserName is the type for the username in the context
type ContextKeyUserName string

const (
	// UserNameKey is the key for the username in the context
	UserNameKey = ContextKeyUserName("username")
)

// DefaultLogFormatter is the default log formatter
func DefaultLogFormatter(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}

	u, err := param.Request.URL.Parse(param.Path)
	// u, err := param.Request.URL.Parse(param.Request.RequestURI)
	if err != nil {
		return "[GIN] invalid request path"
	}

	cuttedPath := u.Path
	// hide /ping prefix
	if len(cuttedPath) > 5 && cuttedPath[:5] == "/ping" {
		cuttedPath = cuttedPath[5:]
	}

	// hide /api prefix
	if len(cuttedPath) > 4 && cuttedPath[:4] == "/api" {
		cuttedPath = cuttedPath[4:]
	}

	if strings.HasPrefix(cuttedPath, "/password/reset/") {
		cuttedPath = "/password/reset/[redacted]"
	}
	queryParams := ""

	// Authentication queries can contain authorization codes, reset secrets and nested URLs.
	// Log the route and status only; never retain query values.
	username := param.Request.Context().Value(UserNameKey)
	if username == nil {
		username = "<nil>"
	}
	if len(username.(string)) > 11 {
		username = username.(string)[:8] + ".."
	}

	return fmt.Sprintf("[GIN] %v %s %3d %s %10v | %-10s |%s %-7s %s %-15s | %s\n%s",
		param.TimeStamp.Format("02.01.2006 - 15:04:05"), // time (%v)
		statusColor, param.StatusCode, resetColor, // status code (%s %3d %s)
		param.Latency,                         // latency (%10v)
		username,                              // username (%-10v)
		methodColor, param.Method, resetColor, // method (%s %-7s %s)
		cuttedPath,         // path (%-15s)
		queryParams,        // query (%#v)
		param.ErrorMessage, // error message (%s)
	)
}
