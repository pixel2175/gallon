package lua

import (
	"os"
	"strconv"

	glua "github.com/yuin/gopher-lua"
)

func optString(t *glua.LTable, key string) (string, bool) {
	s, ok := t.RawGetString(key).(glua.LString)
	return string(s), ok
}

func optInt(t *glua.LTable, key string) (int, bool) {
	n, ok := t.RawGetString(key).(glua.LNumber)
	return int(n), ok
}

func optPositive(L *glua.LState, t *glua.LTable, key string) (int, bool) {
	n, ok := t.RawGetString(key).(glua.LNumber)
	if !ok {
		return 0, false
	}
	if n < 1 {
		L.ArgError(1, key+" must be at least 1")
	}
	return int(n), true
}

func (self *Lua) run(L *glua.LState) int {
	opts := L.CheckTable(1)
	srv := self.State.HttpServer

	srv.Workers = 4
	if n, err := strconv.Atoi(os.Getenv("G_WORKERS")); err == nil && n > 0 {
		srv.Workers = n
	}
	self.State.GC = -1

	if host, ok := optString(opts, "host"); ok {
		srv.Host = host
	}
	if port, ok := optInt(opts, "port"); ok {
		srv.Port = port
	}
	if n, ok := optPositive(L, opts, "workers"); ok {
		srv.Workers = n
	}
	if n, ok := optPositive(L, opts, "gc"); ok {
		self.State.GC = n
	}

	srv.IsRunning = true
	return 0
}
