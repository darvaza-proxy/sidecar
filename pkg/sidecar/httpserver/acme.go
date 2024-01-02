package httpserver

import (
	"net/http"

	"darvaza.org/sidecar/pkg/glob"
)

// AcmeHTTP01Pattern matches ACME-HTTP-01 tokens
var AcmeHTTP01Pattern = glob.MustCompile("/.well-known/acme-challenge{/(*),/**,}", '/')

// AcmeHTTP01Middleware adds middleware to handle HTTP-01 challenges.
// if called without ACME handler or to a path without a valid token
// 404 will be returned automatically instead of passing the request
// to the next handler.
func AcmeHTTP01Middleware(next, acme http.Handler) http.Handler {
	notfound := http.NotFoundHandler()
	if next == nil {
		// end of the line.
		next = notfound
	}

	fn := func(rw http.ResponseWriter, req *http.Request) {
		var h http.Handler

		m, ok := AcmeHTTP01Pattern.Capture(req.URL.Path)

		switch {
		case !ok:
			// not acme-challenge
			h = next
		case m[0] == "":
			// no token
			h = notfound
		case acme == nil:
			// no handler
			h = notfound
		default:
			// handle acme-challenge token
			h = acme
		}
		h.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(fn)
}
