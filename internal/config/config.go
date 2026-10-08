package config

import (
	"errors"
	"net/url"
	"time"

	"github.com/kirillgrachoff/optparf-check/internal/types"
)

type Config struct {
	Telegram *TgConfig     `mapstructure:"telegram"`
	Http     *HttpConfig   `mapstructure:"http"`
	Queries  []QueryConfig `mapstructure:"query"`
	Period   time.Duration `mapstructure:"period"`
}

func (c *Config) Validate() error {
	if c == nil {
		return errors.New("config should be not-null")
	}

	if err := c.Telegram.Validate(); err != nil {
		return err
	}

	return nil
}

type TgConfig struct {
	Token string `mapstructure:"token"`
	Proxy string `mapstructure:"proxy"`
}

func (c *TgConfig) Validate() error {
	if c == nil {
		return nil
	}
	if c.Token == "" {
		return errors.New("tg token is invalid")
	}
	if c.Proxy != "" {
		u, err := url.Parse(c.Proxy)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return errors.New("tg proxy is invalid: expected scheme://host:port")
		}
	}
	return nil
}

type HttpConfig struct {
	QueryPrefix string `mapstructure:"prefix"`
	QuerySuffix string `mapstructure:"suffix"`
}

type QueryConfig = types.Query
