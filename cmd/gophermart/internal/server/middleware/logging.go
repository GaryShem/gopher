package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/logging"
)

type responseData struct {
	status int
	size   int
	body   string
}
type LoggingResponseWriter struct {
	http.ResponseWriter
	data *responseData
}

func (w *LoggingResponseWriter) WriteHeader(statusCode int) {
	w.data.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *LoggingResponseWriter) Write(data []byte) (int, error) {
	if w.data.status == 0 {
		w.data.status = http.StatusOK
		w.ResponseWriter.WriteHeader(http.StatusOK)
	}
	size, err := w.ResponseWriter.Write(data)
	w.data.size += size
	w.data.body = string(data)
	return size, err
}

var _ http.ResponseWriter = &LoggingResponseWriter{}

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
		rData := responseData{}
		next.ServeHTTP(&LoggingResponseWriter{
			ResponseWriter: w,
			data:           &rData}, r)
		logging.Log.Infoln("response code", rData.status)
	})
}
