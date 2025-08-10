package config

import "time"

// Redis 缓存数据库配置
type Redis struct {
	Address      string        `mapstructure:"address" json:"address" yaml:"address"`
	Password     string        `mapstructure:"password" json:"password" yaml:"password"`
	DB           int           `mapstructure:"db" json:"db" yaml:"db"`
	PoolSize     int           `mapstructure:"pool_size" json:"pool_size" yaml:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns" json:"min_idle_conns" yaml:"min_idle_conns"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout" json:"idle_timeout" yaml:"idle_timeout"`
}
