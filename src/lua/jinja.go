package lua

import (
	"gallon/src/jinja"
	"gallon/src/state"

	. "github.com/pixel2175/errs"
	glua "github.com/yuin/gopher-lua"
)

func (self *Lua) ensureJinja() {
	if self.State.Jinja.Loader == nil || self.State.Jinja.Config == nil {
		jinja.Init(self.State.Jinja)
	}
}

func (self *Lua) luaJinjaRenderTemplate(L *glua.LState) int {
	self.ensureJinja()

	name := L.CheckString(1)
	data := L.OptTable(2, L.NewTable())

	result, err := jinja.RenderFile(self.State.Jinja, name, luaTableToMap(data))
	result = TryV(result, err, func(err *error) {
		L.RaiseError("jinja.render_template: %v", *err)
	})

	L.Push(glua.LString(result))
	return 1
}

func (self *Lua) luaJinjaRenderTemplateFromString(L *glua.LState) int {
	self.ensureJinja()

	content := L.CheckString(1)
	data := L.OptTable(2, L.NewTable())

	result, err := jinja.RenderString(self.State.Jinja, content, luaTableToMap(data))

	result = TryV(result, err, func(err *error) {
		L.RaiseError("jinja.render_template_from_string: %v", *err)
	})

	L.Push(glua.LString(result))
	return 1
}

func (self *Lua) luaJinjaSetTemplateDir(L *glua.LState) int {
	dir := L.CheckString(1)

	if self.State.Jinja == nil {
		self.State.Jinja = &state.Jinja{}
	}

	self.State.Jinja.TemplatesDir = dir
	self.State.Jinja.Loader = nil
	self.State.Jinja.Config = nil
	jinja.Init(self.State.Jinja)

	return 0
}
func (self *Lua) registerJinja() {
	jinjaTable := self.Context.NewTable()
	self.Context.SetFuncs(jinjaTable, map[string]glua.LGFunction{
		"render_template":             self.luaJinjaRenderTemplate,
		"render_template_from_string": self.luaJinjaRenderTemplateFromString,
		"set_template_dir":            self.luaJinjaSetTemplateDir,
	})

	self.Gallon.RawSetString("jinja", jinjaTable)
}
