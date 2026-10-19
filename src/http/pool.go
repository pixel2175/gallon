package http

import (
	"net/http"

	"gallon/src/lua"
	"gallon/src/state"

	glua "github.com/yuin/gopher-lua"
)

type worker struct {
	st    *state.State
	reqMT *glua.LTable
}

func (s *Server) initWorkers(n int) {
	s.workers = make(chan *worker, n)

	for range n {
		st := &state.State{
			Server: &state.Server{
				Routes:        map[string]*glua.LFunction{},
				ErrorHandlers: map[string]*glua.LFunction{},
			},
			Args:  s.st.Args,
			Proto: s.st.Proto,
		}
		lua.InitLua(st)
		s.workers <- &worker{
			st:    st,
			reqMT: registerRequestType(st.LuaContext),
		}
	}
}

func (s *Server) acquire(r *http.Request) (*worker, bool) {
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
