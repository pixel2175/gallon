package main

import (
	"gallon/src/run"
	"gallon/src/state"
	"os"

	"github.com/alexflint/go-arg"
)

func main() {
	st := state.Init()
	p := arg.MustParse(st.Args)

	switch {
	case st.Args.Run != nil:
		run.Start(&st)

	default:
		p.WriteHelp(os.Stdout)

	}
}
