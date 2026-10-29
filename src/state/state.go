package state

import (
	"github.com/nikolalohinski/gonja/v2/config"
	"github.com/nikolalohinski/gonja/v2/loaders"
	glua "github.com/yuin/gopher-lua"
)

type Route struct {
	Func     *glua.LFunction
	Body     string
	Status   int
	Headers  map[string]string
	Redirect string
}

type Handlers struct {
	ErrorHandlers map[int]*glua.LFunction
	Routes        map[string]Route
	BeforeRequest *glua.LFunction
	AfterRequest  *glua.LFunction
}

type HttpServer struct {
	IsRunning bool
	Host      string
	Port      int
	Workers   int
	Handlers  Handlers
}
type Jinja struct {
	Loader       loaders.Loader
	Config       *config.Config
	TemplatesDir string
}

type State struct {
	HttpServer *HttpServer
	Args       *Args
	Jinja      *Jinja

	LuaContext *glua.LState
	Proto      *glua.FunctionProto

	GC int
}

func Init() State {
	st := State{
		Args:  &Args{},
		Jinja: &Jinja{},
		HttpServer: &HttpServer{
			IsRunning: false,
			Handlers: Handlers{
				Routes:        make(map[string]Route),
				ErrorHandlers: make(map[int]*glua.LFunction),
			},
		},
	}

	return st
}
