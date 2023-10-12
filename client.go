package bilibili

import (
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"time"
)

type Client struct {
	config *Config

	httpClient *resty.Client
}

func New(cfgs ...Cfg) *Client {
	httpClient := resty.New()

	c := DefaultConfig()
	for _, cfg := range cfgs {
		cfg(c)
	}
	if c.Debug {
		logrus.SetLevel(logrus.DebugLevel)
		httpClient.SetDebug(true)
	}

	logrus.Debugf("create client with config: %+v", c)

	httpClient.
		SetHeader("User-Agent", c.UserAgent).
		SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(20 * time.Second)

	client := &Client{
		config:     c,
		httpClient: httpClient,
	}

	if c.Cookie != nil {
		client.setCookies(c.Cookie)
	}

	return client
}
