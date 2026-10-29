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

type funcs = map[string]glua.LGFunction

func Compile(path string) *glua.FunctionProto {
	var err error
	defer Catch(&err, func(err *error) {
		log.Die("%s", (*err).Error())
	})

	f := TryV(os.Open(path))
	defer f.Close()

	chunk := TryV(parse.Parse(f, path))
	return TryV(glua.Compile(chunk, path))
}

func InitLua(st *state.State) {
	L := glua.NewState(glua.Options{
		RegistrySize:    1024,
		RegistryMaxSize: 20 * 1024,
		CallStackSize:   64,
	})
	st.LuaContext = L

	self := &Lua{State: st, Context: L, Gallon: L.NewTable()}
	self.register()

	L.Push(L.NewFunctionFromProto(st.Proto))
	TryE(L.PCall(0, glua.MultRet, nil))
}

func (self *Lua) register() {
	self.Context.SetFuncs(self.Gallon, funcs{
		"run":            self.run,
		"route":          self.route,
		"before_request": self.beforeRequest,
		"after_request":  self.afterRequest,
		"error_handler":  self.errorHandler,
	})

	self.sub("jinja", funcs{
		"render_template":             self.jinjaRender,
		"render_template_from_string": self.jinjaRenderString,
		"set_template_dir":            self.jinjaSetTemplateDir,
	})

	self.sub("log", funcs{
		"info": self.logInfo,
		"warn": self.logWarn,
		"die":  self.logDie,
	})

	self.Context.SetGlobal("gallon", self.Gallon)
}

func (self *Lua) sub(name string, f funcs) {
	t := self.Context.NewTable()
	self.Context.SetFuncs(t, f)
	self.Gallon.RawSetString(name, t)
}
