package lua

import (
	"strconv"

	glua "github.com/yuin/gopher-lua"
)

func (self *Lua) registerErrorHandler() {
	self.Gallon.RawSetString("error_handler", self.Context.NewFunction(func(L *glua.LState) int {
		code := strconv.Itoa(L.CheckInt(1))
		self.State.Server.ErrorHandlers[code] = L.CheckFunction(2)
		return 0
	}))
}
