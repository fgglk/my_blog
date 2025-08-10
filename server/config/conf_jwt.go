package config

type Jwt struct {
	AccessTokenSecret      string `mapstructure:"access_token_secret" json:"access_token_secret" yaml:"access_token_secret"`
	RefreshTokenSecret     string `mapstructure:"refresh_token_secret" json:"refresh_token_secret" yaml:"refresh_token_secret"`
	AccessTokenExpiryTime  int    `mapstructure:"access_token_expiry_time" json:"access_token_expiry_time" yaml:"access_token_expiry_time"`
	RefreshTokenExpiryTime int    `mapstructure:"refresh_token_expiry_time" json:"refresh_token_expiry_time" yaml:"refresh_token_expiry_time"`
	Issuer                 string `mapstructure:"issuer" json:"issuer" yaml:"issuer"`
}
