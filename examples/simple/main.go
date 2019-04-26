package main

import (
	"flag"
	"fmt"

	"go.coder.com/cli"
)

type cmd struct {
	verbose bool
}

func (c *cmd) Run(fl *flag.FlagSet) {
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

func (c *cmd) RegisterFlags(fl *flag.FlagSet) {
	fl.BoolVar(&c.verbose, "v", false, "sets verbose mode")
}

func main() {
	cli.RunRoot(&cmd{})
}
