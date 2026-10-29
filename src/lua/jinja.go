package lua

import (
	"gallon/src/jinja"
	"gallon/src/state"

	glua "github.com/yuin/gopher-lua"
)

func (self *Lua) jinjaState() *state.Jinja {
	j := self.State.Jinja
	if j == nil {
		j = &state.Jinja{}
		self.State.Jinja = j
	}
	if j.Loader == nil || j.Config == nil {
		jinja.Init(j)
	}
	return j
}

func (self *Lua) jinjaRender(L *glua.LState) int {
	j := self.jinjaState()

	name := L.CheckString(1)
	data := toMap(L.OptTable(2, L.NewTable()))

	result, err := jinja.RenderFile(j, name, data)
	if err != nil {
		L.RaiseError("jinja.render_template: %v", err)
	}

	L.Push(glua.LString(result))
	return 1
}

func (self *Lua) jinjaRenderString(L *glua.LState) int {
	j := self.jinjaState()

	content := L.CheckString(1)
	data := toMap(L.OptTable(2, L.NewTable()))

	result, err := jinja.RenderString(j, content, data)
	if err != nil {
		L.RaiseError("jinja.render_template_from_string: %v", err)
	}

	L.Push(glua.LString(result))
	return 1
}

func (self *Lua) jinjaSetTemplateDir(L *glua.LState) int {
	dir := L.CheckString(1)

	if self.State.Jinja == nil {
		self.State.Jinja = &state.Jinja{}
	}
	j := self.State.Jinja

	j.TemplatesDir = dir
	j.Loader = nil
	j.Config = nil
	jinja.Init(j)

	return 0
}
