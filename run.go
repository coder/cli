package cli

import (
	"os"
	"strings"

	"github.com/spf13/pflag"
)

func appendParent(parent string, add string) string {
	return parent + add + " "
}

// splitArgs tries to split the args between the parent command's flags/args, and the subcommand's
// flags/args. If a subcommand is found, the parent args, subcommand args, and the subcommand will
// be returned. If a subcommand isn't found, the args will be returned as is, the subcommand args
// will be empty, and the subcommand will be nil.
func splitArgs(subCmds []Command, args []string) (cmdArgs, subArgs []string, subCmd Command) {
	for i, arg := range args {
		if strings.HasPrefix(arg, "-") {
			continue
		}

		for _, subCommand := range subCmds {
			if subCommand.Spec().Name == arg {
				return args[:i], args[i+1:], subCommand
			}
		}
	}

	return args, []string{}, nil
}

// Run sets up flags, helps, and executes the command with the provided
// arguments.
//
// parents is the list of parent commands.
// E.g the parent for `sail run hello` would be `sail run`.
//
// Use RunRoot if this package is managing the entire CLI.
func Run(cmd Command, args []string, parent string) {
	name := parent + cmd.Spec().Name
	fl := pflag.NewFlagSet(name, pflag.ContinueOnError)
	// Ensure pflag library doesn't print usage for us automatically,
	// we'll override this below.
	fl.Usage = func() {}

	if fc, ok := cmd.(FlaggedCommand); ok {
		fc.RegisterFlags(fl)
	}

	if cmd.Spec().RawArgs {
		// Use `--` to return immediately when parsing the flags.
		args = append([]string{"--"}, args...)
	}

	var (
		cmdArgs, subArgs []string
		subCmd           Command
	)
	pc, isParentCmd := cmd.(ParentCommand)
	if isParentCmd {
		cmdArgs, subArgs, subCmd = splitArgs(pc.Subcommands(), args)
		if subCmd != nil {
			args = cmdArgs
		}
	}

	err := fl.Parse(args)
	// Reassign the usage now that we've parsed the args
	// so that we can render it manually.
	fl.Usage = func() {
		renderHelp(name, cmd, fl, os.Stderr)
	}
	if err != nil {
		fl.Usage()
		os.Exit(2)
	}

	// Route to subcommand.
	if isParentCmd && subCmd != nil {
		Run(
			subCmd, subArgs,
			appendParent(parent, cmd.Spec().Name),
		)
		return
	}

	cmd.Run(fl)
}

// RunRoot calls Run with the process's arguments.
func RunRoot(cmd Command) {
	Run(cmd, os.Args[1:], "")
}
