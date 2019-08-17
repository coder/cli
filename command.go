package cli

import (
	"flag"
	"strings"
)

// Command describes a command or subcommand.
type Command interface {
	// Configure adds metadata about the command.
	// The command must set the Name, Usage, and Desc fields.
	Configure(s *Config)

	// Run invokes the command's main routine with parsed flags.
	Run(s *Config)
}

// Config describes a command's configuration.
type Config struct {
	// Name is the name of the command.
	// It should be the leaf name of the entire command. E.g `run` for the full
	// command `sail run`.
	// Required.
	Name string

	// Usage is the command's usage string.
	// E.g "[flags] <path>"
	// Required.
	Usage string

	// Desc is the description of the command.
	// The first line is used as an abbreviated description.
	// Desc should not list flags.
	// Required.
	Desc string

	// RawArgs indicates that flags should not be parsed, and they should be deferred
	// to the command.
	RawArgs bool

	// Hidden indicates that this command should not show up in it's parent's
	// subcommand help.
	Hidden bool

	// Flags registers all of the commands flags.
	Flags *flag.FlagSet

	// Subcommands a list of children commands.
	// The parent command may pass itself into its children in order to easily share dependencies and state.
	Subcommands []Command
}

// ShortDesc returns the first line of Desc.
func (c Config) ShortDesc() string {
	return strings.Split(c.Desc, "\n")[0]
}

// AddSubcommand is a helper which adds a subcommand.
func (c Config) AddSubcommand(sc Command) {
	c.Subcommands = append(c.Subcommands, sc)
}

func (c Config) isParent() bool {
	return len(c.Subcommands) > 0
}

func (c Config) isValid() bool {
	return c.Name != "" && c.Usage != "" && c.Desc != ""
}
