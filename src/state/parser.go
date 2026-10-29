package state

type CmdRun struct {
	MainPoint string `arg:"positional" default:"main.lua" help:"main lua file"`
}

type Args struct {
	Run *CmdRun `arg:"subcommand:run" help:"Run web server"`
}
