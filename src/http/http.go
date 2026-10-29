package http

import (
	"fmt"
	"gallon/src/log"
	"gallon/src/state"
	"net/http"
	"runtime/debug"
	"strings"

	. "github.com/pixel2175/errs"
)

type Server struct {
	st      *state.State
	workers chan *worker
}

func Start(st *state.State) {
	if !st.HttpServer.IsRunning {
		log.Die("HTTP server is not running; call gallon.run()")
	}
	s := &Server{st: st}
	s.initWorkers()
	mux := http.NewServeMux()

	for path := range s.st.HttpServer.Handlers.Routes {
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

	addr := fmt.Sprintf("%s:%d", s.st.HttpServer.Host, s.st.HttpServer.Port)
	log.Info("Server Starts on: %s", addr)
	TryE(http.ListenAndServe(addr, mux))

}
