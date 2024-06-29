package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/logging"
)

func LogBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			logging.Log.Infoln("no basic auth provided")
		} else {
			logging.Log.Infof("basic auth provided: %s:%s\n", username, password)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logging.Log.Infoln(string(body))
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		next.ServeHTTP(w, r)
		logging.Log.Infoln("response code", w.Header().Get("status"))
	})
}
