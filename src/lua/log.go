package lua

import (
	"gallon/src/log"

	glua "github.com/yuin/gopher-lua"
)

func (self *Lua) luaLogDie(L *glua.LState) int {
	msg := L.CheckString(1)
	log.Die(msg)
	return 0
}

func (self *Lua) luaLogWarn(L *glua.LState) int {
	msg := L.CheckString(1)
	log.Warn(msg)
	return 0
}

func (self *Lua) luaLogInfo(L *glua.LState) int {
	msg := L.CheckString(1)
	log.Info("%s", msg)
	return 0
}

func (self *Lua) registerLog() {
	logTable := self.Context.NewTable()
	self.Context.SetFuncs(logTable, map[string]glua.LGFunction{
		"info": self.luaLogInfo,
		"warn": self.luaLogWarn,
		"die":  self.luaLogDie,
	})

	self.Gallon.RawSetString("log", logTable)
}
