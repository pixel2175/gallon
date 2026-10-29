package http

import (
	"errors"
	"net/http"

	"gallon/src/lua"

	. "github.com/pixel2175/errs"
	glua "github.com/yuin/gopher-lua"
)

var errNoHandler = errors.New("route has no handler")

type Response struct {
	Body     string
	Status   int
	Headers  map[string]string
	Redirect string
}

func run(w *worker, fn *glua.LFunction, r *http.Request, hooks bool) (resp Response, err error) {
	defer Catch(&err)

	if fn == nil {
		return resp, errNoHandler
	}

	L, srv := w.luaContext, w.server
	req, res := lua.RequestToTable(L, r), L.CreateTable(0, 0)

	if hooks && srv.Handlers.BeforeRequest != nil {
		TryE(invoke(L, srv.Handlers.BeforeRequest, 0, req, res))
	}

	resp = TryV(handle(L, fn, req, res))
	if hooks && srv.Handlers.AfterRequest != nil {
		TryE(invoke(L, srv.Handlers.AfterRequest, 0, req, res))
	}
	return resp, nil
}

func invoke(L *glua.LState, fn *glua.LFunction, nret int, req, res glua.LValue) error {
	L.Push(fn)
	L.Push(req)
	L.Push(res)
	return L.PCall(2, nret, nil)
}

func handle(L *glua.LState, fn *glua.LFunction, req, res glua.LValue) (resp Response, err error) {
	defer Catch(&err)

	top := L.GetTop()
	defer L.SetTop(top)
	TryE(invoke(L, fn, 3, req, res))
	body, status, headers := L.Get(-3), L.Get(-2), L.Get(-1)

	resp = Response{Status: http.StatusOK}

	if body != glua.LNil {
		resp.Body = body.String()
	}
	if n, ok := status.(glua.LNumber); ok {
		resp.Status = int(n)
	}
	if t, ok := headers.(*glua.LTable); ok {
		resp.Headers = make(map[string]string)
		t.ForEach(func(k, v glua.LValue) {
			resp.Headers[k.String()] = v.String()
		})
	}
	return resp, nil
}
