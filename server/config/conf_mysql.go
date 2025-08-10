package config

import (
	"time"
)

type Mysql struct {
	Host         string        `json:"host" yaml:"host" mapstructure:"host"`
	Port         int           `json:"port" yaml:"port" mapstructure:"port"`
	Username     string        `json:"username" yaml:"username" mapstructure:"username"`
	Password     string        `json:"password" yaml:"password" mapstructure:"password"`
	DBName       string        `json:"db_name" yaml:"db_name" mapstructure:"db_name"`
	Config       string        `json:"config" yaml:"config" mapstructure:"config"`
	MaxIdleConns int           `json:"max_idle_conns" yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
	MaxOpenConns int           `json:"max_open_conns" yaml:"max_open_conns" mapstructure:"max_open_conns"`
	LogMode      string        `json:"log_mode" yaml:"log_mode" mapstructure:"log_mode"`
	MaxLifeTime  time.Duration `json:"max_life_time" yaml:"max_life_time" mapstructure:"max_life_time"`
	MaxIdleTime  time.Duration `json:"conn_max_idle_time" yaml:"conn_max_idle_time" mapstructure:"conn_max_idle_time"`
}
