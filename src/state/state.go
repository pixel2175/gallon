package state

import glua "github.com/yuin/gopher-lua"

type CmdRun struct {
	MainPoint string `arg:"positional" default:"main.lua" help:"main lua file"`
}

type Args struct {
	Run *CmdRun `arg:"subcommand:run" help:"Run web server"`
}

type Server struct {
	IsRunning bool
	Host      string
	Port      int
	Workers   int

	ErrorHandlers map[string]*glua.LFunction
	Routes        map[string]*glua.LFunction
	BeforeRequest *glua.LFunction
	AfterRequest  *glua.LFunction
}

type State struct {
	Server     *Server
	Args       *Args
	LuaContext *glua.LState
	Proto      *glua.FunctionProto
	GC         int
}

func Init() State {
	st := State{
		Args: &Args{},
		Server: &Server{
			IsRunning:     false,
			Routes:        make(map[string]*glua.LFunction),
			ErrorHandlers: make(map[string]*glua.LFunction),
		},
	}

	return st
}
