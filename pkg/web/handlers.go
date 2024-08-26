package web

import (
	"fmt"
	"net/http"
	"net/url"

	httpSwagger "github.com/swaggo/http-swagger"
)

// RedirectHandlerFunc - служит для быстрого создания редиректов
func RedirectHandlerFunc(code int, to string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, to, code)
	}
}

// SwaggerHandler - служит для вызова страницы с документацией
func SwaggerHandler(baseURL *url.URL, swaggerFilePath string) http.Handler {
	return httpSwagger.Handler(
		httpSwagger.URL(baseURL.Path+swaggerFilePath),
		httpSwagger.BeforeScript(`const UrlMutatorPlugin = (system) => ({
			rootInjects: {
			  setScheme: (scheme) => {
				const jsonSpec = system.getState().toJSON().spec.json;
				const schemes = Array.isArray(scheme) ? scheme : [scheme];
				const newJsonSpec = Object.assign({}, jsonSpec, { schemes })
				return system.specActions.updateJsonSpec(newJsonSpec);
			  },
			  setHost: (host) => {
				const jsonSpec = system.getState().toJSON().spec.json;
				const newJsonSpec = Object.assign({}, jsonSpec, { host })
				return system.specActions.updateJsonSpec(newJsonSpec);
			  },
			  setBasePath: (basePath) => {
				const jsonSpec = system.getState().toJSON().spec.json;
				const newJsonSpec = Object.assign({}, jsonSpec, { basePath })
				return system.specActions.updateJsonSpec(newJsonSpec);
			  }
			}
		});`),
		httpSwagger.Plugins([]string{"UrlMutatorPlugin"}),
		httpSwagger.UIConfig(map[string]string{
			"onComplete": fmt.Sprintf(`() => {
				window.ui.setScheme('%s');
				window.ui.setHost('%s');
				window.ui.setBasePath('%s');
			}`, baseURL.Scheme, baseURL.Host, baseURL.Path),
		}),
	)
}

// JSON_NotFoundHandlerFunc - создает стандартную заглушку
// для ситуации, когда не удалось найти handler для обработки
// запроса
func JSON_NotFoundHandlerFunc(w http.ResponseWriter, r *http.Request) {
	NewResponse[error](
		fmt.Errorf("method %s not found for path %s", r.Method, r.URL.Path)).
		Write(http.StatusNotFound, w)
}

// JSON_MethodNotAllowedHandlerFunc - создает стандартную заглушку
// для ситуации, когда не удалось найти handler для обработки
// запроса
func JSON_MethodNotAllowedHandlerFunc(w http.ResponseWriter, r *http.Request) {
	NewResponse[error](
		fmt.Errorf("method %s not allowed for path %s", r.Method, r.URL.Path)).
		Write(http.StatusMethodNotAllowed, w)
}
