package lua

import (
	glua "github.com/yuin/gopher-lua"
)

func (self *Lua) registerRoutes() {

	self.Gallon.RawSetString("route", self.Context.NewFunction(func(L *glua.LState) int {
		path := L.CheckString(1)
		self.State.Server.Routes[path] = L.CheckFunction(2)
		return 0
	}))
}
