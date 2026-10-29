package http

import (
	"fmt"
	"gallon/src/log"
	"io"
	"net/http"
	"os"

	. "github.com/pixel2175/errs"
)

func (s *Server) route(path string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var err error
		defer Catch(&err, func(err *error) {
			log.Warn(log.Title("ERROR"), "%s", (*err).Error())
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			*err = nil
		})
		worker, ok := s.acquire(r)
		if !ok {
			return
		}
		defer s.release(worker)
		route := worker.server.Handlers.Routes[path]
		resp, err := run(worker, route.Func, r, true)
		TryV(resp, err, func(err error) error {
			fmt.Fprintln(os.Stderr, err)
			resp = errorPage(worker, r, http.StatusInternalServerError)
			return nil
		})

		Write(resp, rw)

	}
}

func (s *Server) notFound(rw http.ResponseWriter, r *http.Request) {
	worker, ok := s.acquire(r)
	if !ok {
		return
	}
	defer s.release(worker)

	resp := errorPage(worker, r, http.StatusNotFound)
	Write(resp, rw)
}

func Write(resp Response, w http.ResponseWriter) {
	for k, v := range resp.Headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(resp.Status)
	io.WriteString(w, resp.Body)
}

func errorPage(worker *worker, r *http.Request, code int) Response {
	resp := Response{Body: http.StatusText(code)}
	if fn := worker.server.Handlers.ErrorHandlers[code]; fn != nil {
		res, err := run(worker, fn, r, true)
		switch {
		case err == nil:
			resp = res
		case code != http.StatusInternalServerError:
			fmt.Fprintln(os.Stderr, err)
			return errorPage(worker, r, http.StatusInternalServerError)
		default:
			fmt.Fprintln(os.Stderr, err)
		}
	}
	resp.Status = code
	return resp
}
