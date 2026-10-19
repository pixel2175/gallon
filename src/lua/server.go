package lua

import (
	"os"
	"strconv"

	glua "github.com/yuin/gopher-lua"
)

func (self *Lua) registerRun() {
	self.Gallon.RawSetString("run", self.Context.NewFunction(func(L *glua.LState) int {
		opts := L.CheckTable(1)
		srv := self.State.Server

		srv.Workers = 4
		self.State.GC = -1
		if n, err := strconv.Atoi(os.Getenv("G_WORKERS")); err == nil && n > 0 {
			srv.Workers = n
		}

		if host, ok := opts.RawGetString("host").(glua.LString); ok {
			srv.Host = string(host)
		}
		if port, ok := opts.RawGetString("port").(glua.LNumber); ok {
			srv.Port = int(port)
		}
		if n, ok := opts.RawGetString("workers").(glua.LNumber); ok {
			if n < 1 {
				L.ArgError(1, "workers must be at least 1")
			}
			srv.Workers = int(n)
		}
		if n, ok := opts.RawGetString("gc").(glua.LNumber); ok {
			if n < 1 {
				L.ArgError(1, "gc must be at least 1")
			}
			self.State.GC = int(n)
		}

		srv.IsRunning = true
		return 0
	}))
}
