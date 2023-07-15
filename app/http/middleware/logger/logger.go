// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"

	"web-tpl/app/core/config"
	"web-tpl/app/utils/env"
)

type consoleColorModeValue int

const (
	autoColor consoleColorModeValue = iota
	disableColor
	forceColor
)

const logTimeTpl = "2006-01-02T15:04:05.000Z07:00"
const LogID = "X-Log-Id"

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

var consoleColorMode = autoColor

// LoggerConfig defines the config for Logger middleware.
type LoggerConfig struct {
	// Optional. Default value is gin.defaultLogFormatter
	Formatter LogFormatter

	// Output is a writer where logs are written.
	// Optional. Default value is gin.DefaultWriter.
	Output io.Writer

	// SkipPaths is an url path array which logs are not written.
	// Optional.
	SkipPaths []string
}

// LogFormatter gives the signature of the formatter function passed to LoggerWithFormatter
type LogFormatter func(params LogFormatterParams) string

// LogFormatterParams is the structure any formatter will be handed when time to log comes
type LogFormatterParams struct {
	LocalIP  string
	Env      string // 环境
	Hostname string
	Format   string // 文件输出的格式
	LogID    string

	Request *http.Request

	// TimeStamp shows the time after the server returns a response.
	TimeStamp time.Time
	// StatusCode is HTTP response code.
	StatusCode int
	// Latency is how much time the server cost to process a certain request.
	Latency time.Duration
	// ClientIP equals Context's ClientIP method.
	ClientIP string
	// Method is the HTTP method given to the request.
	Method string
	// Path is a path the client requests.
	Path string
	// ErrorMessage is set if error has occurred in processing the request.
	ErrorMessage string
	// BodySize is the size of the Response Body
	BodySize int
	// Keys are the keys set on the request's context.
	Keys map[string]any
}

// StatusCodeColor is the ANSI color for appropriately logging http status code to a terminal.
func (p *LogFormatterParams) StatusCodeColor() string {
	code := p.StatusCode

	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return yellow
	default:
		return red
	}
}

// MethodColor is the ANSI color for appropriately logging http method to a terminal.
func (p *LogFormatterParams) MethodColor() string {
	method := p.Method

	switch method {
	case http.MethodGet:
		return blue
	case http.MethodPost:
		return cyan
	case http.MethodPut:
		return yellow
	case http.MethodDelete:
		return red
	case http.MethodPatch:
		return green
	case http.MethodHead:
		return magenta
	case http.MethodOptions:
		return white
	default:
		return reset
	}
}

// ResetColor resets all escape attributes.
func (p *LogFormatterParams) ResetColor() string {
	return reset
}

// IsOutputColor indicates whether can colors be outputted to the log.
func (p *LogFormatterParams) IsOutputColor() bool {
	return consoleColorMode == forceColor || consoleColorMode == autoColor
}

// defaultLogFormatter is the default log format function Logger middleware uses.
var defaultLogFormatter = func(param LogFormatterParams) string {

	// 支持text, json
	if param.Format == "json" {
		var rel = map[string]any{
			"ip":         param.ClientIP,
			"log_id":     param.LogID,
			"time":       param.TimeStamp.Format(logTimeTpl),
			"latency":    param.Latency.Milliseconds(),
			"method":     param.Method,
			"path":       param.Path,
			"query":      param.Request.URL.RawQuery,
			"status":     param.StatusCode,
			"error":      param.ErrorMessage,
			"size":       param.BodySize,
			"local_ip":   param.LocalIP,
			"env":        param.Env,
			"hostname":   param.Hostname,
			"user_agent": param.Request.UserAgent(),
			"referer":    param.Request.Referer(),
		}

		relJson, _ := json.Marshal(rel)
		return string(relJson) + "\n"
	}

	// text
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	return fmt.Sprintf(
		"[%s]\t%s\t%s\t%d\t%s%s%s\t%s\t%s\t%s%d%s\t%s\t%d\t%s\t%s\t%s\t%s\n",
		param.ClientIP,
		param.LogID,
		param.TimeStamp.Format(logTimeTpl),
		param.Latency.Milliseconds(),
		methodColor, param.Method, resetColor,
		param.Path,
		param.Request.URL.RawQuery,
		statusColor, param.StatusCode, resetColor,
		param.ErrorMessage,
		param.BodySize,
		param.Env,
		param.Hostname,
		param.Request.UserAgent(),
		param.Request.Referer(),
	)
}

// DisableConsoleColor disables color output in the console.
func DisableConsoleColor() {
	consoleColorMode = disableColor
}

// ForceConsoleColor force color output in the console.
func ForceConsoleColor() {
	consoleColorMode = forceColor
}

var output io.Writer = os.Stdout

func New(conf config.WebServerLog, prjEnv string, homeDir string) gin.HandlerFunc {
	// 解决当前输出是stdout还是file
	switch conf.Output {
	case "file":
		var err error
		output, err = loadLogFile(conf, homeDir)
		if err != nil {
			panic(err)
		}
	default:
		output = os.Stdout
	}

	localIP := env.LocalIP()
	hostname := env.Hostname()

	formatter := defaultLogFormatter

	return func(c *gin.Context) {
		var logId = c.GetHeader(LogID)
		if logId == "" {
			logId = xid.New().String()
			c.Request.Header.Add(LogID, logId)
		}

		// Start timer
		start := time.Now()
		path := c.Request.URL.Path

		if conf.LogIDShowHeader {
			c.Header(LogID, logId)
		}

		// Process request
		c.Next()

		var skip map[string]struct{}
		if length := len(conf.SkipPaths); length > 0 {
			skip = make(map[string]struct{}, length)
			for _, p := range conf.SkipPaths {
				skip[p] = struct{}{}
			}
		}

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			param := LogFormatterParams{
				Request: c.Request,
				Keys:    c.Keys,
			}

			// Stop timer
			param.Format = conf.LogFormat
			param.LocalIP = localIP
			param.Hostname = hostname
			param.TimeStamp = time.Now()
			param.Latency = param.TimeStamp.Sub(start)
			param.LogID = logId
			param.ClientIP = c.ClientIP()
			param.Method = c.Request.Method
			param.StatusCode = c.Writer.Status()
			param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
			param.BodySize = c.Writer.Size()
			param.Env = prjEnv
			param.Path = path

			_, _ = fmt.Fprint(output, formatter(param))
		}
	}
}

func loadLogFile(conf config.WebServerLog, homeDir string) (io.Writer, error) {
	logPath := "logs/access.log"
	if conf.LogPath != "" {
		logPath = conf.LogPath
	}

	// 判断logPath是相对路径还是绝对路径
	if !filepath.IsAbs(logPath) {
		logPath = homeDir + "/" + logPath
	}

	// 检测这个文件是否存在，如果不存在呢我们就创建这个文件
	f, e := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if e != nil {
		return nil, e
	}

	return f, nil
}
