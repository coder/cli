# cli

A minimal, command-oriented CLI package.

[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/go.coder.com/cli)

## Features

- Very small, simple API.
- Support for POSIX flags.
- Only external dependency is [spf13/pflag](https://github.com/spf13/pflag).
- Subcommands.
- Auto-generated help.

## Install

```bash
go get -u go.coder.com/cli
```

## Examples

See `examples/` for more.

### Simple CLI
```go
package main

import (
    "flag"
    "fmt"

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
    fl.BoolVar(&c.verbose, "v", false, "sets verbose mode")
}

func main() {
    cli.RunRoot(&cmd{})
}

```
renders a help like

```
Usage: simple-example [flags]

This is a simple example of the cli package.

simple-example flags:
	-v	sets verbose mode	(false)
```

### Subcommands

```go
package main

import (
    "flag"
    "fmt"

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
```

renders a help like

```
Usage: subcommand [flags]

This is a simple example of subcommands.

Commands:
	sub	This is a simple subcommand.
```