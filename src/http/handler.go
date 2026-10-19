package http

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	. "github.com/pixel2175/errs"
)

func (s *Server) route(path string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		w, ok := s.acquire(r)
		if !ok {
			return
		}

		resp, err := run(w, w.st.Server.Routes[path], r, true)
		TryV(resp, err, func(err error) error {
			fmt.Fprintln(os.Stderr, err)
			resp = errorPage(w, r, http.StatusInternalServerError)
			return nil
		})

		s.release(w)
		resp.Write(rw)
	}
}

func (s *Server) notFound(rw http.ResponseWriter, r *http.Request) {
	w, ok := s.acquire(r)
	if !ok {
		return
	}

	resp := errorPage(w, r, http.StatusNotFound)
	s.release(w)
	resp.Write(rw)
}

func errorPage(w *worker, r *http.Request, code int) Response {
	resp := Response{Body: http.StatusText(code)}

	if fn := w.st.Server.ErrorHandlers[strconv.Itoa(code)]; fn != nil {
		if got, err := run(w, fn, r, false); err == nil {
			resp = got
		}
	}

	resp.Status = code
	return resp
}
