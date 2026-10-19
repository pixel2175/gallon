package http

import (
	"io"
	"net/http"
)

type Response struct {
	Body    string
	Status  int
	Headers map[string]string
}

func (resp Response) Write(w http.ResponseWriter) {
	for k, v := range resp.Headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(resp.Status)
	io.WriteString(w, resp.Body)
}
