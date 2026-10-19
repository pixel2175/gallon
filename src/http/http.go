package http

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"gallon/src/state"
	. "github.com/pixel2175/errs"
)

type Server struct {
	st      *state.State
	workers chan *worker
}

func Start(st *state.State) {
	if !st.Server.IsRunning {
		return
	}

	s := &Server{st: st}
	s.initWorkers(s.st.Server.Workers)

	mux := http.NewServeMux()

	for path := range s.st.Server.Routes {
		if path == "/" || strings.HasSuffix(path, " /") {
			mux.HandleFunc(path+"{$}", s.route(path))
		} else {
			mux.HandleFunc(path, s.route(path))
		}
	}

	mux.HandleFunc("/", s.notFound)

	if s.st.GC > 0 {
		debug.SetGCPercent(st.GC)
	}

	addr := fmt.Sprintf("%s:%d", s.st.Server.Host, s.st.Server.Port)
	TryE(http.ListenAndServe(addr, mux))
}
