package admin

import (
	"Noctobot/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Noctember/gocto"
	"github.com/Noctember/gocto/helpers"
	"github.com/dop251/goja"
	"github.com/jonas747/discordgo"
	"github.com/mattn/anko/env"
	"github.com/mattn/anko/vm"
	"github.com/starlight-go/starlight/convert"
	"go.starlark.net/starlark"
	"strings"
)

func OwnerEval(ctx *gocto.CommandContext) {
	vm := goja.New()

	vm.Set("ctx", ctx)
	vm.Set("bot", ctx.Bot)
	vm.Set("session", ctx.Session)
	vm.Set("print", fmt.Print)

	value, err := vm.RunString(ctx.JoinedArgs())
	if err != nil {
		ctx.CodeBlock("js", "%s", err.Error())
		return
	}

	if value.String() != "" {
		ctx.CodeBlock("js", value.String())
	}

}

func OwnerAnko(ctx *gocto.CommandContext) {
	e := env.NewEnv()
	stdout := &bytes.Buffer{}

	e.Define("println", func(args ...interface{}) {
		stdout.WriteString(fmt.Sprintln(args...))
	})
	e.Define("NewCommand", gocto.NewCommand)
	e.Define("PermToText", helpers.GetPermissionsText)

	e.Define("ctx", ctx)
	e.Define("bot", ctx.Bot)
	e.Define("session", ctx.Session)
	e.Define("ksoft", utils.Ksoft)
	e.Define("hypixel", utils.Hypixel)
	e.Define("crab", "🦀")
	e.Define("SnowflakeTimestamp", utils.SnowflakeTimestamp)
	e.Define("Marshal", json.Marshal)
	e.Define("printb", func(args ...interface{}) {
		stdout.WriteString(fmt.Sprintf("%s\n", args...))
	})

	e.Define("aep", discordgo.EndpointAPI)

	value, err := vm.Execute(e, nil, ctx.JoinedArgs())
	if err != nil {
		ctx.CodeBlock("go", "%v\n", err)
		return
	}
	if stdout.String() == "<nil>" {
		return
	}

	ctx.CodeBlock("go", "%s%+v", stdout.String(), value)
}

func OwnerStarlark(ctx *gocto.CommandContext) {
	code := strings.Trim(strings.Trim(ctx.Message.Content[len(ctx.Prefix)+len(ctx.InvokedName):], " "), "```")
	var stdout []string

	thread := &starlark.Thread{Name: "eval", Print: func(_ *starlark.Thread, msg string) {
		stdout = append(stdout, msg)
	}}

	dict, err := convert.MakeStringDict(map[string]interface{}{
		"ctx":     ctx,
		"bot":     ctx.Bot,
		"session": ctx.Session,
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	_, err = starlark.ExecFile(thread, "eval", code, dict)
	if err != nil {
		ctx.Error(err)
		return
	}
	if err != nil {
		if trace, ok := err.(*starlark.EvalError); ok {
			ctx.CodeBlock("py", trace.Backtrace())
		} else {
			ctx.CodeBlock("py", err.Error())
		}

		return
	}

	if len(stdout) < 1 {
		return
	}

	ctx.CodeBlock("py", "%+v", strings.Join(stdout, "\n"))
}
