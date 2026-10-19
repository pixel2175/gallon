package lua

import (
	"gallon/src/log"
	"gallon/src/state"
	"os"

	. "github.com/pixel2175/errs"
	glua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
)

type Lua struct {
	State   *state.State
	Context *glua.LState
	Gallon  *glua.LTable
}

func Compile(path string) *glua.FunctionProto {
	var err error
	defer Catch(&err, func(err *error) {
		log.Die("%s", (*err).Error())
	})

	f := TryV(os.Open(path))

	defer f.Close()

	chunk := TryV(parse.Parse(f, path))
	p := TryV(glua.Compile(chunk, path))
	return p
}

func InitLua(st *state.State) {
	st.LuaContext = glua.NewState(glua.Options{
		RegistrySize:    1024,
		RegistryMaxSize: 1024 * 20,
		CallStackSize:   64,
	})

	self := Lua{
		State:   st,
		Context: st.LuaContext,
	}

	self.registerGallon()

	L := st.LuaContext
	L.Push(L.NewFunctionFromProto(st.Proto))
	TryE(L.PCall(0, glua.MultRet, nil))

}

func (self *Lua) registerGallon() {
	self.Gallon = self.State.LuaContext.NewTable()

	self.registerRun()
	self.registerMiddleware()
	self.registerRoutes()
	self.registerErrorHandler()
	self.registerLog()

	self.State.LuaContext.SetGlobal("gallon", self.Gallon)
}
