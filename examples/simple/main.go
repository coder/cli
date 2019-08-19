package main

import (
	"fmt"

	"github.com/spf13/pflag"
	"go.coder.com/cli"
)

type cmd struct {
	verbose bool
}

func (c *cmd) Run(fl *pflag.FlagSet) {
	if c.verbose {
		fmt.Println("verbose enabled")
	}
	fmt.Println("we run")
}

func (c *cmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "simple-example",
		Usage: "[flags]",
		Desc:  `This is a simple example of the cli package.`,
	}
}

func (c *cmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.BoolVarP(&c.verbose, "verbose", "v", false, "sets verbose mode")
}

func main() {
	cli.RunRoot(&cmd{})
}
