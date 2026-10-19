package http

import (
	"gallon/src/log"
	"net/http"

	. "github.com/pixel2175/errs"
	glua "github.com/yuin/gopher-lua"
)

func run(w *worker, fn *glua.LFunction, r *http.Request, hooks bool) (resp Response, err error) {
	defer Catch(&err)
	L := w.st.LuaContext
	req, res := newRequest(L, w.reqMT, r), L.NewTable()

	if hooks {
		TryE(call(L, w.st.Server.BeforeRequest, req, res))
	}

	resp = TryV(callResponse(L, fn, req, res))

	if hooks {
		TryE(call(L, w.st.Server.AfterRequest, req, res))
	}
	return resp, nil
}

func call(L *glua.LState, fn *glua.LFunction, req glua.LValue, res *glua.LTable) error {
	if fn == nil {
		return nil
	}
	return L.CallByParam(glua.P{Fn: fn, NRet: 0, Protect: true}, req, res)
}

func callResponse(L *glua.LState, fn *glua.LFunction, req glua.LValue, res *glua.LTable) (resp Response, err error) {
	defer Catch(&err)

	if fn == nil {
		log.Info("OS")
	}

	TryE(L.CallByParam(glua.P{Fn: fn, NRet: 3, Protect: true}, req, res))
	body, status, headers := L.Get(-3), L.Get(-2), L.Get(-1)
	L.Pop(3)

	resp = Response{Status: http.StatusOK}

	switch b := body.(type) {
	case glua.LString:
		resp.Body = string(b)
	default:
		if body != glua.LNil {
			resp.Body = body.String()
		}
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
