package logger

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	// LoggerConfig is the config for the logger
	LoggerConfig = gin.LoggerConfig{
		SkipPaths: []string{"/internal"},
		Formatter: DefaultLogFormatter,
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
		return "[GIN] " + err.Error()
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

	queryParams := ""

	m, err := url.ParseQuery(u.RawQuery)
	if err == nil {
		for k, v := range m {
			value := strings.Join(v, ",")
			if param.IsOutputColor() {
				if value == "true" {
					value = "\033[97;42mtrue\033[0m"
				}
				if value == "false" {
					value = "\033[97;41mfalse\033[0m"
				}
			}
			queryParams = queryParams + fmt.Sprintf("%s=%s  ", k, value)
		}
	}
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
