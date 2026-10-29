package jinja

import (
	"errors"
	"fmt"
	"strings"

	"gallon/src/state"

	. "github.com/pixel2175/errs"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"
	"github.com/nikolalohinski/gonja/v2/parser"
)

func Init(j *state.Jinja) {
	if j.TemplatesDir == "" {
		j.TemplatesDir = "templates/"
	}

	j.Loader = TryV(loaders.NewFileSystemLoader(j.TemplatesDir))
	j.Config = gonja.DefaultConfig
}

func RenderString(j *state.Jinja, content string, data map[string]any) (result string, err error) {
	defer Catch(&err)

	loader := TryV(loaders.NewShiftedLoader(
		"root",
		strings.NewReader(content),
		j.Loader,
	))

	tpl := TryV(exec.NewTemplate(
		"root",
		j.Config,
		loader,
		gonja.DefaultEnvironment,
	))

	return render(tpl, data)
}

func RenderFile(j *state.Jinja, name string, data map[string]any) (result string, err error) {
	defer Catch(&err)

	tpl := TryV(exec.NewTemplate(
		name,
		j.Config,
		j.Loader,
		gonja.DefaultEnvironment,
	))

	return render(tpl, data)
}

func render(tpl *exec.Template, data map[string]any) (result string, err error) {
	ctx := exec.NewContext(data)

	v, e := tpl.ExecuteToString(ctx)
	result = TryV(v, e, func(err *error) {
		if syntaxErr, ok := errors.AsType[*parser.SyntaxError](*err); ok {
			*err = fmt.Errorf(
				"Jinja syntax error:\n- Line=%d, Column=%d\n- %s",
				syntaxErr.Line,
				syntaxErr.Column,
				syntaxErr.Message,
			)
		}
	})
	return
}
