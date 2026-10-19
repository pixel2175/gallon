package http

import (
	"bytes"
	"io"
	"net/http"
	"regexp"

	glua "github.com/yuin/gopher-lua"
)

var wildcard = regexp.MustCompile(`\{(\w+)(?:\.\.\.)?\}`)

type request struct {
	r     *http.Request
	body  []byte
	read  bool
	cache map[string]glua.LValue
}

func registerRequestType(L *glua.LState) *glua.LTable {
	mt := L.NewTypeMetatable("request")
	L.SetField(mt, "__index", L.NewFunction(requestIndex))
	return mt
}

func newRequest(L *glua.LState, mt *glua.LTable, r *http.Request) glua.LValue {
	ud := L.NewUserData()
	ud.Value = &request{r: r}
	L.SetMetatable(ud, mt)
	return ud
}

func requestIndex(L *glua.LState) int {
	q := L.CheckUserData(1).Value.(*request)
	key := L.CheckString(2)

	v, ok := q.cache[key]
	if !ok {
		v = q.field(L, key)
		if q.cache == nil {
			q.cache = map[string]glua.LValue{}
		}
		q.cache[key] = v
	}
	L.Push(v)
	return 1
}

func (q *request) field(L *glua.LState, key string) glua.LValue {
	r := q.r
	switch key {
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
		return glua.LString(q.readBody())
	case "headers":
		return lists(L, r.Header)
	case "query":
		return lists(L, r.URL.Query())
	case "form":
		q.readBody()
		r.ParseForm()
		return lists(L, r.Form)
	case "cookies":
		return cookies(L, r)
	case "path_value":
		return pathValue(L, r)
	}
	return glua.LNil
}

func (q *request) readBody() []byte {
	if !q.read {
		q.body, _ = io.ReadAll(q.r.Body)
		q.r.Body = io.NopCloser(bytes.NewReader(q.body))
		q.read = true
	}
	return q.body
}

func lists(L *glua.LState, m map[string][]string) *glua.LTable {
	t := L.NewTable()
	for k, values := range m {
		list := L.NewTable()
		for _, v := range values {
			list.Append(glua.LString(v))
		}
		t.RawSetString(k, list)
	}
	return t
}

func cookies(L *glua.LState, r *http.Request) *glua.LTable {
	t := L.NewTable()
	for _, c := range r.Cookies() {
		t.RawSetString(c.Name, glua.LString(c.Value))
	}
	return t
}

func pathValue(L *glua.LState, r *http.Request) *glua.LFunction {
	return L.NewFunction(func(L *glua.LState) int {
		if L.GetTop() > 0 {
			L.Push(glua.LString(r.PathValue(L.CheckString(1))))
			return 1
		}

		all := L.NewTable()
		for _, m := range wildcard.FindAllStringSubmatch(r.Pattern, -1) {
			all.RawSetString(m[1], glua.LString(r.PathValue(m[1])))
		}
		L.Push(all)
		return 1
	})
}
