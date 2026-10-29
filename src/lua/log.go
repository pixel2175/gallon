package lua

import (
	"gallon/src/log"

	glua "github.com/yuin/gopher-lua"
)

func (self *Lua) logInfo(L *glua.LState) int {
	log.Info("%s", L.CheckString(1))
	return 0
}

func (self *Lua) logWarn(L *glua.LState) int {
	log.Warn(L.CheckString(1))
	return 0
}

func (self *Lua) logDie(L *glua.LState) int {
	log.Die(L.CheckString(1))
	return 0
}
