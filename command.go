package cli

import (
	"flag"
	"strings"
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
}

// ShortDesc returns the first line of Desc.
func (c CommandSpec) ShortDesc() string {
	return strings.Split(c.Desc, "\n")[0]
}

// Command describes a command or subcommand.
type Command interface {
	// Spec returns metadata about the command.
	Spec() CommandSpec
	// Run invokes the command's main routine with parsed flags.
	Run(fl *flag.FlagSet)
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
	RegisterFlags(fl *flag.FlagSet)
}
