package config

type Email struct {
	Host      string `mapstructure:"host" json:"host" yaml:"host"`
	Port      int    `mapstructure:"port" json:"port" yaml:"port"`
	From      string `mapstructure:"from" json:"from" yaml:"from"`
	NNickname string `mapstructure:"nickname" json:"nickname" yaml:"nickname"`
	Secret    string `mapstructure:"secret" json:"secret" yaml:"secret"`
	IsSsl     bool   `mapstructure:"is_ssl" json:"is_ssl" yaml:"is_ssl"`
}
