package config

import (
	"time"
)

type Captcha struct {
	Height     int           `mapstructure:"height" json:"height" yaml:"height"`
	Width      int           `mapstructure:"width" json:"width" yaml:"width"`
	Length     int           `mapstructure:"length" json:"length" yaml:"length"`
	MaxSkew    float64       `mapstructure:"max_skew" json:"max_skew" yaml:"max_skew"`
	DotCount   int           `mapstructure:"dot_count" json:"dot_count" yaml:"dot_count"`
	Expiration time.Duration `mapstructure:"expiration" json:"expiration" yaml:"expiration"`
}
