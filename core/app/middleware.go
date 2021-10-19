package app

import (
	"net/http"
	"strings"

	"github.com/urfave/negroni"
)

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
func TrimTrailingSlash(allSlashes bool) Middleware {
	if !allSlashes {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
				next.ServeHTTP(rw, r)
			})
		}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			for len(r.URL.Path) > 0 && r.URL.Path[len(r.URL.Path)-1] == '/' {
				r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
			}
			next.ServeHTTP(rw, r)
		})
	}
}
