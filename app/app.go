package app

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"web-tpl/app/core/log"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"web-tpl/app/core/config"
	"web-tpl/app/core/db"
	rds "web-tpl/app/core/redis"
)

var Config config.Model

func Init(prjHome string) error {
	// 找到两个配置文件路径
	return Config.LoadConfig(prjHome)

	// viper
}

func DBW(keys ...string) *gorm.DB {
	k := "default"
	if len(keys) > 0 {
		k = keys[0]
	}

	conf, ok := Config.DB[k]
	if !ok {
		panic(fmt.Sprintf("db config %s not found", k))
	}
	cacheKey := fmt.Sprintf("%s_write", k)

	return db.Load(conf.Write, conf.Log, cacheKey, Config.Env, Config.HomeDir)
}

func DBR(keys ...string) *gorm.DB {
	k := "default"
	if len(keys) > 0 {
		k = keys[0]
	}

	conf, ok := Config.DB[k]
	if !ok {
		panic(fmt.Sprintf("db config %s not found", k))
	}

	cacheKey := fmt.Sprintf("%s_read", k)

	return db.Load(conf.Read, conf.Log, cacheKey, Config.Env, Config.HomeDir)
}

func Log(keys ...string) *logrus.Entry {
	return log.Load(Config.HomeDir, Config.Log, Config.Env)
}

func RedisR(key ...string) *redis.Client {
	return rds.Load("read", &Config, key...)
}

func RedisW(key ...string) *redis.Client {
	return rds.Load("write", &Config, key...)
}
