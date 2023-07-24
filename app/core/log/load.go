package log

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"sync"
	"web-tpl/app/core/config"
	"web-tpl/app/utils/env"
)

const logTimeTpl = "2006-01-02T15:04:05.000Z07:00"

var log *logrus.Entry
var logLocker sync.RWMutex

func Load(homeDir string, logConf config.Log, currentEnv string) *logrus.Entry {
	logLocker.RLock()
	if log != nil {
		logLocker.RUnlock()
		return log
	}
	logLocker.RUnlock()

	// A, B, C
	logLocker.Lock()
	defer logLocker.Unlock()

	// 二次判断
	if log != nil {
		return log
	}

	logNew := logrus.New()

	logNew.SetReportCaller(true)

	// 设置日志level
	setLogLevel(logNew, logConf.Level)

	// 设置日志格式， json 或者 text {
	setLogFormat(logNew, logConf.LogFormat, currentEnv)

	// 设置日式输出
	setLogOutput(logNew, logConf, homeDir)

	// 基础字段预设，比如项目名、环境、env、local_ip、hostname、idc
	l := presetFields(logNew, currentEnv)
	log = l

	return l
}

func presetFields(logger *logrus.Logger, currentEnv string) *logrus.Entry {
	return logger.WithFields(logrus.Fields{
		"env":      currentEnv,
		"local_ip": env.LocalIP(),
		"hostname": env.Hostname(),
		"module":   filepath.Base(os.Args[0]),
	})
}

func setLogOutput(logger *logrus.Logger, logConf config.Log, homeDir string) {
	if logConf.Output == "file" {
		f, e := loadLogFile(logConf, homeDir)
		if e != nil {
			panic(e)
		}
		logger.SetOutput(f)
	} else {
		logger.SetOutput(os.Stdout)
	}
}

func setLogFormat(logger *logrus.Logger, format, currentEnv string) {
	if format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: logTimeTpl,
		})
	} else {
		// 如果非dev环境禁用掉color
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: logTimeTpl,
			DisableColors:   currentEnv != "dev",
		})
	}
}

func loadLogFile(conf config.Log, homeDir string) (io.Writer, error) {
	logPath := "logs/app.log"
	if conf.Name != "" {
		logPath = conf.Name
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

func setLogLevel(logger *logrus.Logger, level string) {
	switch level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}
}
