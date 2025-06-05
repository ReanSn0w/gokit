package web

import "net/url"

func NewRestConfig(revision string, baseURL *url.URL) *RestConfig {
	return &RestConfig{
		Revision: revision,
		BaseURL:  baseURL,

		DocsPath:   "/docs",
		StaticPath: "/static",

		DocsFilePath: "/static/swagger.yaml",
		StaticDir:    "static",
	}
}

type RestConfig struct {
	Revision string   // Версия приложения для отображения в документации (Например v1.0.0)
	BaseURL  *url.URL // BaseURL вервиса (для локального запуска обычно http://localhost:8080)

	DocsPath   string // Path по которому будет находится Swagger handler (по умолчанию: /docs)
	StaticPath string // Path по которому будет находится путь со статическими файлами (по умолчанию: /static)

	DocsFilePath string // путь к файлу документации на севере (по умолчанию: /static/swagger.yaml)
	StaticDir    string // путь к директории статических файлов (по умолчанию: ./static)
}

func (r *RestConfig) SetDocsPath(val string) *RestConfig {
	r.DocsPath = val
	return r
}

func (r *RestConfig) SetStaticPath(val string) *RestConfig {
	r.StaticPath = val
	return r
}

func (r *RestConfig) SetDocsFilePath(val string) *RestConfig {
	r.DocsFilePath = val
	return r
}

func (r *RestConfig) SetStaticDir(val string) *RestConfig {
	r.StaticDir = val
	return r
}
