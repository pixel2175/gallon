package http

import (
	"gallon/src/lua"
	"gallon/src/state"
	"net/http"

	glua "github.com/yuin/gopher-lua"
)

type worker struct {
	server     *state.HttpServer
	luaContext *glua.LState
}

func (s *Server) initWorkers() {
	n := s.st.HttpServer.Workers
	s.workers = make(chan *worker, n)

	for range n {
		st := state.Init()
		st.Proto = s.st.Proto
		lua.InitLua(&st)

		s.workers <- &worker{
			server:     st.HttpServer,
			luaContext: st.LuaContext,
		}
	}
}

func (s *Server) acquire(r *http.Request) (*worker, bool) {
	select {
	case w := <-s.workers:
		return w, true
	default:
	}
	select {
	case w := <-s.workers:
		return w, true
	case <-r.Context().Done():
		return nil, false
	}
}

func (s *Server) release(w *worker) {
	s.workers <- w
}
