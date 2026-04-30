package web

import (
	"net/http"
	"strings"
)

// swaggerPageTMPL is an HTML page template for the Scalar API Reference UI.
// It contains two placeholders that are replaced at construction time:
//   - {{.Title}}  — the page <title> shown in the browser tab
//   - {{.DocURL}} — the URL of the OpenAPI/Swagger specification file
const (
	swaggerPageTMPL = `
<!DOCTYPE html>
<html lang="ru">

<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>{{.Title}}</title>
</head>

<body>
  <div id="scalar-api-reference"></div>
  <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  <script>
    Scalar.createApiReference('#scalar-api-reference', {
      url: '{{.DocURL}}',
      showDeveloperTools: "never",
    })
  </script>
</body>

</html>
`
)

// NewScalarHandler returns an HTTP handler that serves an HTML page embedding
// the Scalar API Reference UI (loaded from CDN).
//
// Parameters:
//   - title  — the page title rendered in the browser tab (<title>)
//   - docURL — the URL of the OpenAPI/Swagger specification that Scalar will load
//
// The handler always responds with HTTP 200 and Content-Type text/html.
func NewScalarHandler(title string, docURL string) http.HandlerFunc {
	pageTMPL := strings.NewReplacer(
		"{{.Title}}", title,
		"{{.DocURL}}", docURL,
	).Replace(swaggerPageTMPL)

	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(pageTMPL))
	}
}
