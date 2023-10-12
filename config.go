package bilibili

import (
	"github.com/mcuadros/go-defaults"
)

type Config struct {
	Debug    bool `default:"false"`
	Validate bool `default:"false"`

	Cookie    *Cookie
	UserAgent string `default:"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36"`
}

func DefaultConfig() *Config {
	config := new(Config)
	defaults.SetDefaults(config)
	return config
}

type Cfg func(config *Config)

func WithUserAgent(userAgent string) Cfg {
	return func(config *Config) {
		config.UserAgent = userAgent
	}
}

func WithCookie(cookie *Cookie) Cfg {
	return func(config *Config) {
		config.Cookie = cookie
	}
}

func WithDebug(debug bool) Cfg {
	return func(config *Config) {
		config.Debug = debug
	}
}
