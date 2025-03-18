package config

import (
	"net/http"
	"os"
	"time"

	"github.com/go-pkgz/lgr"
	"golang.org/x/net/proxy"
)

type HTTPClientConfig struct {
	HTTPClient struct {
		Timeout int `long:"timeout" env:"TIMEOUT" default:"60" description:"http response timeout (seconds)"`

		Proxy struct {
			Enabled  bool   `long:"enabled" env:"ENABLED" description:"enable http socks proxy"`
			Host     string `long:"host" env:"HOST" description:"proxy host"`
			Login    string `long:"login" env:"login" description:"proxy login"`
			Password string `long:"password" env:"PASSWORD" description:"proxy password"`
		} `group:"proxy" namespace:"proxy" env-namespace:"PROXY"`
	} `group:"http-client" namespace:"http-client" env-namespace:"HTTP_CLIENT"`
}

func (config *HTTPClientConfig) HTTPClientMustCreate() *http.Client {
	client, err := config.HTTPClientCreate()
	if err != nil {
		lgr.Default().Logf("[ERROR] create http client err: $v", err)
		os.Exit(2)
	}
	return client
}

func (config *HTTPClientConfig) HTTPClientCreate() (*http.Client, error) {
	client := &http.Client{
		Timeout: time.Second * time.Duration(config.HTTPClient.Timeout),
	}

	if config.HTTPClient.Proxy.Enabled {
		var auth *proxy.Auth
		if config.HTTPClient.Proxy.Login != "" {
			auth = &proxy.Auth{
				User:     config.HTTPClient.Proxy.Login,
				Password: config.HTTPClient.Proxy.Password,
			}
		}

		dialer, err := proxy.SOCKS5("tcp", config.HTTPClient.Proxy.Host, auth, nil)
		if err != nil {
			return nil, err
		}

		client.Transport = &http.Transport{Dial: dialer.Dial}
	}

	return client, nil
}
