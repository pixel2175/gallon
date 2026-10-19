package run

import (
	"gallon/src/http"
	"gallon/src/log"
	"gallon/src/lua"
	"gallon/src/state"

	. "github.com/pixel2175/errs"
)

func Start(state *state.State) (err error) {
	defer Catch(&err, func(err *error) {
		log.Die("%s", (*err).Error())
	})

	state.Proto = lua.Compile(state.Args.Run.MainPoint)
	lua.InitLua(state)
	http.Start(state)
	defer state.LuaContext.Close()

	return nil
}
