package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"web-tpl/app/core/config"
)

const (
	dbTimeout      = 5000 // ms
	dbWriteTimeout = 5000 // ms
	dbReadTimeout  = 5000 // ms
)

var dbInstance = make(map[string]*gorm.DB)
var dbLocker sync.RWMutex

func Load(conf config.DBItemConf, confLog config.DBLog, key string, env string, homeDir string) *gorm.DB {
	dbLocker.RLock()
	db, ok := dbInstance[key]
	if ok {
		dbLocker.RUnlock()
		return db
	}
	dbLocker.RUnlock()

	dbLocker.Lock()
	defer dbLocker.Unlock()

	// double check
	if _, exist := dbInstance[key]; exist {
		return dbInstance[key]
	}

	dbInstance[key] = getDBInstance(conf, confLog, env, homeDir)

	return dbInstance[key]
}

func getDBInstance(conf config.DBItemConf, confLog config.DBLog, env string, homeDir string) *gorm.DB {
	timeout := dbTimeout
	if conf.TimeOut > 0 {
		timeout = conf.TimeOut
	}

	writeTimeout := dbWriteTimeout
	if conf.WriteTimeOut > 0 {
		writeTimeout = conf.WriteTimeOut
	}

	readTimeout := dbReadTimeout
	if conf.ReadTimeOut > 0 {
		readTimeout = conf.ReadTimeOut
	}

	dsnConf := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local&timeout=%dms&writeTimeout=%dms&readTimeout=%dms",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database,
		conf.Charset,
		timeout,
		writeTimeout,
		readTimeout,
	)

	var l logger.LogLevel
	var dbLogger *log.Logger
	var gLogger logger.Interface

	if confLog.Enable {
		switch confLog.Level {
		case "silent":
			l = logger.Silent
		case "error":
			l = logger.Error
		case "info":
			l = logger.Info
		default:
			l = logger.Warn
		}

		if confLog.Type == "file" {
			logPath := confLog.Path
			if !filepath.IsAbs(confLog.Path) {
				logPath = homeDir + "/" + confLog.Path
			}

			f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				panic(err)
			}
			dbLogger = log.New(f, "", log.LstdFlags)
		} else if confLog.Type == "stdout" {
			dbLogger = log.New(os.Stdout, "", log.LstdFlags)
		}

		gLogger = New(dbLogger, logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  l,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		}, env)
	}

	db, err := gorm.Open(mysql.Open(dsnConf), &gorm.Config{
		Logger: gLogger,
	})

	if err != nil {
		panic(err)
	}

	// 数据库池配置
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConns)

	return db
}
