package rw

import "net/http"

type ResponseWriter struct {
	http.ResponseWriter
	Code int
	Err  error
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.Code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *ResponseWriter) WriteError(err error) {
	w.Err = err
}
