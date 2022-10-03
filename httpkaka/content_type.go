package httpkaka

import (
	"net/http"
	"strings"
)

func setContentTypeJSON(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		handler.ServeHTTP(w, r)
	}
}

func hasExtension(path string, extension string) bool {
	return strings.HasSuffix(path, extension) || strings.HasSuffix(path, extension+"/")
}

func setContentTypeFileFormat(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		value := "text/html"

		if hasExtension(r.URL.Path, ".css") {
			value = "text/css"
		}

		if hasExtension(r.URL.Path, ".js") {
			value = "text/javascript"
		}

		w.Header().Add("Content-Type", value)
		handler.ServeHTTP(w, r)
	})
}
