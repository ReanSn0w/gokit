package config

import (
	"net/url"
	"os"

	"github.com/go-pkgz/lgr"
)

type HTTPServerConfig struct {
	HTTPServer struct {
		Port int    `long:"port" env:"PORT" default:"8080" description:"http server port"`
		URL  string `long:"url" env:"URL" default:"http://localhost:8080" descripiton:"http server base url"`
	} `group:"http-server" namespace:"http-server" env-namespace:"HTTP_SERVER"`
}

func (h HTTPServerConfig) BaseURL() *url.URL {
	url, err := url.Parse(h.HTTPServer.URL)
	if err != nil {
		lgr.Default().Logf("[ERROR] parse base url err: %v", err)
		os.Exit(2)
	}

	return url
}
