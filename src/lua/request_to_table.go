package lua

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	glua "github.com/yuin/gopher-lua"
)

type reqCtx struct {
	r      *http.Request
	body   []byte
	read   bool
	pathFn *glua.LFunction
}

func RequestToTable(L *glua.LState, r *http.Request) *glua.LTable {
	mt, ok := L.GetTypeMetatable("gallon.req").(*glua.LTable)
	if !ok {
		mt = newReqMeta(L)
	}
	ctx := mt.RawGetString("__ctx").(*glua.LUserData).Value.(*reqCtx)
	ctx.r, ctx.body, ctx.read = r, nil, false

	t := L.CreateTable(0, 16)
	L.SetMetatable(t, mt)
	return t
}

func newReqMeta(L *glua.LState) *glua.LTable {
	ctx := &reqCtx{}
	ctx.pathFn = L.NewFunction(ctx.pathValue)

	mt := L.NewTypeMetatable("gallon.req")
	ud := L.NewUserData()
	ud.Value = ctx
	mt.RawSetString("__ctx", ud)
	mt.RawSetString("__index", L.NewFunction(func(L *glua.LState) int {
		t, key := L.CheckTable(1), L.CheckString(2)
		v := ctx.field(L, key)
		if v != glua.LNil {
			t.RawSetString(key, v)
		}
		L.Push(v)
		return 1
	}))
	return mt
}

func (c *reqCtx) readBody() {
	if c.read {
		return
	}
	c.read = true
	if c.r.Body != nil && c.r.Body != http.NoBody {
		c.body, _ = io.ReadAll(c.r.Body)
		c.r.Body = io.NopCloser(bytes.NewReader(c.body))
	}
}

func (c *reqCtx) field(L *glua.LState, k string) glua.LValue {
	r := c.r
	switch k {
	case "method":
		return glua.LString(r.Method)
	case "path":
		return glua.LString(r.URL.Path)
	case "url":
		return glua.LString(r.URL.String())
	case "pattern":
		return glua.LString(r.Pattern)
	case "host":
		return glua.LString(r.Host)
	case "proto":
		return glua.LString(r.Proto)
	case "remote_addr":
		return glua.LString(r.RemoteAddr)
	case "raw_query":
		return glua.LString(r.URL.RawQuery)
	case "content_length":
		return glua.LNumber(r.ContentLength)
	case "body":
		c.readBody()
		return glua.LString(c.body)
	case "headers":
		return list(L, r.Header)
	case "query":
		return list(L, r.URL.Query())
	case "form":
		c.readBody()
		r.ParseForm()
		if c.body != nil {
			r.Body = io.NopCloser(bytes.NewReader(c.body))
		}
		return list(L, r.Form)
	case "cookies":
		cs := r.Cookies()
		t := L.CreateTable(0, len(cs))
		for _, ck := range cs {
			t.RawSetString(ck.Name, glua.LString(ck.Value))
		}
		return t
	case "path_value":
		return c.pathFn
	}
	return glua.LNil
}

func (c *reqCtx) pathValue(L *glua.LState) int {
	r := c.r
	if L.GetTop() > 0 {
		L.Push(glua.LString(r.PathValue(L.CheckString(1))))
		return 1
	}
	all := L.CreateTable(0, 0)
	p := r.Pattern
	// i hate my life
	for {
		i := strings.IndexByte(p, '{')
		if i < 0 {
			break
		}
		j := strings.IndexByte(p[i:], '}')
		if j < 0 {
			break
		}
		if name := strings.TrimSuffix(p[i+1:i+j], "..."); name != "$" {
			all.RawSetString(name, glua.LString(r.PathValue(name)))
		}
		p = p[i+j+1:]
	}
	L.Push(all)
	return 1
}

func list(L *glua.LState, m map[string][]string) *glua.LTable {
	t := L.CreateTable(0, len(m))
	for k, vs := range m {
		l := L.CreateTable(len(vs), 0)
		for _, v := range vs {
			l.Append(glua.LString(v))
		}
		t.RawSetString(k, l)
	}
	return t
}
