package lua

import (
	"gallon/src/state"

	glua "github.com/yuin/gopher-lua"
)

func (self *Lua) route(L *glua.LState) int {
	self.State.HttpServer.Handlers.Routes[L.CheckString(1)] = state.Route{Func: L.CheckFunction(2)}
	return 0
}

func (self *Lua) beforeRequest(L *glua.LState) int {
	self.State.HttpServer.Handlers.BeforeRequest = L.CheckFunction(1)
	return 0
}

func (self *Lua) afterRequest(L *glua.LState) int {
	self.State.HttpServer.Handlers.AfterRequest = L.CheckFunction(1)
	return 0
}

func (self *Lua) errorHandler(L *glua.LState) int {
	code := L.CheckInt(1)
	self.State.HttpServer.Handlers.ErrorHandlers[code] = L.CheckFunction(2)
	return 0
}
