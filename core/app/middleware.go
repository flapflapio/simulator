package app

import (
	"net/http"

	"github.com/urfave/negroni"
)

type Middleware func(http.Handler) http.Handler

func LoggerAndRecovery(h http.Handler) http.Handler {
	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.Use(negroni.NewRecovery())
	n.UseHandler(h)
	return n
}

// Returns a middleware that trims trailing slashes. If `allSlashes` is `true`,
// the middlware trims all trailing slashes, otherwise the middleware trims only
// a single trailing slash
func TrimTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		r.URL.Path = Trim(r.URL.Path)
		next.ServeHTTP(rw, r)
	})
}

func Trim(url string) string {
	for l := len(url); l > 1 && url[l-2:] == "//"; l = len(url) {
		url = url[:l-1]
	}
	return url
}
