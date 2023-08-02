package config

import "time"

// RedisItem 配置
type RedisItem struct {
	Addr         string        `yaml:"addr"`
	Password     string        `yaml:"password"`
	PoolSize     int           `yaml:"pool_size"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
	MinIdleConns int           `yaml:"min_idle_conns"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	Retries      int           `yaml:"retries"`
	DB           int           `yaml:"db"`
}

type Redis struct {
	Write RedisItem `yaml:"write"`
	Read  RedisItem `yaml:"read"`
}
