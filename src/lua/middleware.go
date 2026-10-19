package lua

import (
	glua "github.com/yuin/gopher-lua"
)

func (self *Lua) registerMiddleware() {
	self.Gallon.RawSetString("before_request", self.Context.NewFunction(func(L *glua.LState) int {
		self.State.Server.BeforeRequest = L.CheckFunction(1)
		return 0
	}))

	self.Gallon.RawSetString("after_request", self.Context.NewFunction(func(L *glua.LState) int {
		self.State.Server.AfterRequest = L.CheckFunction(1)
		return 0
	}))
}
