package redis

import (
	"fmt"
	"sync"
	"web-tpl/app/core/config"

	"github.com/go-redis/redis/v8"
)

var redisClientSyncMap map[string]*redis.Client
var redisClientMutex sync.RWMutex

// Load redis instance. flag: read or write
func Load(flag string, conf *config.Model, key ...string) *redis.Client {
	var k = "default"
	if len(key) != 0 {
		k = key[0] // use default key name
	}

	sgKey := flag + k
	redisClientMutex.RLock()
	if c, ok := redisClientSyncMap[sgKey]; ok {
		redisClientMutex.RUnlock()
		return c
	}
	redisClientMutex.RUnlock()

	// double-checked locking
	redisClientMutex.Lock()
	defer redisClientMutex.Unlock()
	if c, ok := redisClientSyncMap[sgKey]; ok {
		return c
	}

	rdsConf, ok := conf.Redis[k]
	if !ok {
		panic(fmt.Sprintf("redis %s config not exist, please check your yaml config!", key[0]))
	}

	c := initRedis(flag == "write", &rdsConf)
	redisClientSyncMap[sgKey] = c

	return c
}

// valid addr empty panic
func validParams(addr string) {
	if addr == "" {
		panic("redis can not allow empty addr")
	}
}

func initRedis(isWrite bool, conf *config.Redis) *redis.Client {
	//var rdsConf *redis.Options
	var rdsConf config.RedisItem
	if isWrite {
		rdsConf = conf.Write
	} else {
		rdsConf = conf.Read
	}
	validParams(rdsConf.Addr)

	return redis.NewClient(&redis.Options{
		Addr:         rdsConf.Addr,
		PoolSize:     rdsConf.PoolSize,
		Password:     rdsConf.Password,
		IdleTimeout:  rdsConf.IdleTimeout,
		ReadTimeout:  rdsConf.IdleTimeout,
		WriteTimeout: rdsConf.WriteTimeout,
		MaxRetries:   rdsConf.Retries,
		MinIdleConns: rdsConf.MinIdleConns,
		DB:           rdsConf.DB,
	})
}
