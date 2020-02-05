package apply

import (
	"flag"
	"fmt"
	"os"

	"github.com/m3db/m3/src/cmd/tools/m3ctl/main/errors"
	"github.com/m3db/m3/src/cmd/tools/m3ctl/main/globalopts"
	"github.com/m3db/m3/src/x/config/configflag"
)

type applyVals struct {
	yamlFlag *configflag.FlagStringSlice
}
type Context struct {
	vals       *applyVals
	handlers   applyHandlers
	GlobalOpts globalopts.GlobalOpts
	Flags      *flag.FlagSet
}
type applyHandlers struct {
	handle func(*applyVals, globalopts.GlobalOpts)
}

func InitializeFlags() Context {
	return _setupFlags(
		&applyVals{
			yamlFlag: &configflag.FlagStringSlice{},
		},
		applyHandlers{
			handle: doApply,
		})
}

func _setupFlags(finalArgs *applyVals, handlers applyHandlers) Context {
	applyFlags := flag.NewFlagSet("apply", flag.ContinueOnError)
	applyFlags.Var(finalArgs.yamlFlag, "f", "Path to yaml.")
	applyFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, "help msg here for apply\n")
		applyFlags.PrintDefaults()
	}
	return Context{
		vals:     finalArgs,
		Flags:    applyFlags,
		handlers: handlers,
	}
}

func (ctx Context) PopParseDispatch(cli []string) error {
	if len(cli) < 1 {
		ctx.Flags.Usage()
		return &errors.FlagsError{}
	}
	inArgs := cli[1:]
	if err := ctx.Flags.Parse(inArgs); err != nil {
		ctx.Flags.Usage()
		return err
	}
	if ctx.Flags.NFlag() != 1 {
		ctx.Flags.Usage()
		return &errors.FlagsError{}
	}
	if err := dispatcher(ctx); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return err
	}
	return nil
}

// no dispatching here
// there are no subcommands
func dispatcher(ctx Context) error {
	ctx.handlers.handle(ctx.vals, ctx.GlobalOpts)
	return nil
}
