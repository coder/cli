package cli

import (
	"flag"
	"os"
)

func appendParent(parent string, add string) string {
	if parent == "" {
		return add + " "
	}
	return parent + add + " "
}

// Run executes sets up flags, helps, and executes the command with the provided
// arguments.
//
// parents is the list of parent commands.
// E.g the parent for `sail run hello` would be `sail run`.
//
// Use RunRoot if this package is managing the entire CLI.
func Run(cmd Command, args []string, parent string) {
	fl := flag.NewFlagSet(parent+""+cmd.Spec().Name, flag.ExitOnError)

	if fc, ok := cmd.(FlaggedCommand); ok {
		fc.RegisterFlags(fl)
	}

	fl.Usage = func() {
		renderHelp(cmd, fl, os.Stderr)
	}
	_ = fl.Parse(args)

	subcommandArg := fl.Arg(0)

	// Route to subcommand.
	if pc, ok := cmd.(ParentCommand); ok && subcommandArg != "" {
		for _, subcommand := range pc.Subcommands() {
			if subcommand.Spec().Name != subcommandArg {
				continue
			}

			Run(
				subcommand, fl.Args()[1:],
				appendParent(parent, cmd.Spec().Name),
			)
			return
		}
	}

	cmd.Run(fl)
}

// RunRoot calls Run with the process's arguments.
func RunRoot(cmd Command) {
	Run(cmd, os.Args[1:], "")
}
