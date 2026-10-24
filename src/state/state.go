package state

import (
	"github.com/nikolalohinski/gonja/v2/config"
	"github.com/nikolalohinski/gonja/v2/loaders"
	glua "github.com/yuin/gopher-lua"
)

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
type Jinja struct {
	Loader loaders.Loader
	Config *config.Config
	TemplatesDir string
}

type State struct {
	Server     *Server
	Args       *Args
	Jinja      *Jinja

	LuaContext *glua.LState
	Proto      *glua.FunctionProto

	GC         int
}

func Init() State {
	st := State{
		Args: &Args{},
		Jinja: &Jinja{},
		Server: &Server{
			IsRunning:     false,
			Routes:        make(map[string]*glua.LFunction),
			ErrorHandlers: make(map[string]*glua.LFunction),
		},
	}

	return st
}
