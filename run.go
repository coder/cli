package cli

import (
	"flag"
	"math/rand"
	"os"
	"strconv"
)

func appendParent(parent string, add string) string {
	if parent == "" {
		return add + " "
	}
	return parent + add + " "
}

// commandConfig gets a configuration from a command.
func commandConfig(cmd Command) *Config {
	c := &Config{
		Flags: flag.NewFlagSet(strconv.FormatInt(rand.Int63(), 10), flag.ExitOnError),
	}
	cmd.Configure(c)
	c.Flags.Usage = func() {
		renderHelp(c, os.Stderr)
	}
	return c
}

// Run sets up flags, helps, and executes the command with the provided
// arguments.
//
// parents is the list of parent commands.
// E.g the parent for `sail run hello` would be `sail run`.
//
// Use RunRoot if this package is managing the entire CLI.
func Run(cmd Command, args []string, parent string) {
	c := commandConfig(cmd)

	if c.RawArgs {
		// Use `--` to return immediately when parsing the flags.
		args = append([]string{"--"}, args...)
	}

	_ = c.Flags.Parse(args)
	subcommandArg := c.Flags.Arg(0)

	// Route to subcommand.
	if c.isParent() && subcommandArg != "" {
		for _, subcommand := range c.Subcommands {
			sc := commandConfig(subcommand)
			if sc.Name != subcommandArg {
				continue
			}

			Run(
				subcommand, c.Flags.Args()[1:],
				appendParent(parent, sc.Name),
			)
			return
		}
	}

	cmd.Run(c)
}

// RunRoot calls Run with the process's arguments.
func RunRoot(cmd Command) {
	Run(cmd, os.Args[1:], "")
}
