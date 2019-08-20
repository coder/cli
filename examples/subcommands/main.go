package main

import (
	"fmt"

	"github.com/spf13/pflag"
	"go.coder.com/cli"
)

type subcmd struct {
}

func (c *subcmd) Run(fl *pflag.FlagSet) {
	fmt.Println("subcommand invoked")
}

func (c *subcmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "sub",
		Usage: "",
		Desc:  `This is a simple subcommand.`,
	}
}

type cmd struct {
}

func (c *cmd) Run(fl *pflag.FlagSet) {
	// This root command has no default action, so print the help.
	fl.Usage()
}

func (c *cmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "subcommand",
		Usage: "[flags]",
		Desc:  `This is a simple example of subcommands.`,
	}
}

func (c *cmd) Subcommands() []cli.Command {
	return []cli.Command{
		&subcmd{},
	}
}

func main() {
	cli.RunRoot(&cmd{})
}
