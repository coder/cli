package cli

import (
	"strings"

	"github.com/spf13/pflag"
)

// CommandSpec describes a Command's usage.
//
// It should not list flags.
type CommandSpec struct {
	// Name is the name of the command.
	// It should be the leaf name of the entire command. E.g `run` for the full
	// command `sail run`.
	Name string
	// Usage is the command's usage string.
	// E.g "[flags] <path>"
	Usage string
	// Desc is the description of the command.
	// The first line is used as an abbreviated description.
	Desc string
	// RawArgs indicates that flags should not be parsed, and they should be deferred
	// to the command.
	RawArgs bool
	// Hidden indicates that this command should not show up in it's parent's
	// subcommand help.
	Hidden bool
	// Aliases contains a list of alternative names that can be used for a particular command.
	Aliases []string
}

// ShortDesc returns the first line of Desc.
func (c CommandSpec) ShortDesc() string {
	return strings.Split(c.Desc, "\n")[0]
}

// HasAliases evaluates whether particular command has any alternative names.
func (c CommandSpec) HasAliases() bool {
	return len(c.Aliases) > 0
}

// Command describes a command or subcommand.
type Command interface {
	// Spec returns metadata about the command.
	Spec() CommandSpec
	// Run invokes the command's main routine with parsed flags.
	Run(fl *pflag.FlagSet)
}

// ParentCommand is an optional interface for commands that have subcommands.
//
// A ParentCommand may pass itself into children as it creates them in order to
// pass high-level configuration and state.
type ParentCommand interface {
	Subcommands() []Command
}

// FlaggedCommand is an optional interface for commands that have flags.
type FlaggedCommand interface {
	// RegisterFlags lets the command register flags which be sent to Handle.
	RegisterFlags(fl *pflag.FlagSet)
}
