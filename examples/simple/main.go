package main

import (
	"fmt"

	"go.coder.com/cli"
)

type cmd struct {
	verbose bool
}

func (c *cmd) Run(cf *cli.Config) {
	if c.verbose {
		fmt.Println("verbose enabled")
	}
	fmt.Println("we run")
}

func (c *cmd) Configure(cf *cli.Config) {
	cf.Name = "simple-example"
	cf.Usage = "[flags]"
	cf.Desc = `This is a simple example of the cli package.`

	cf.Flags.BoolVar(&c.verbose, "v", false, "sets verbose mode")
}

func main() {
	cli.RunRoot(&cmd{})
}
