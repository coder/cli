# cli

A minimal, command-oriented CLI package.

[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/go.coder.com/cli)

## Features

- Very small, simple API.
- Uses standard `flag` package as much as possible.
- No external dependencies.
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

func (c *subcmd) Run(fl *flag.FlagSet) {
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

func (c *cmd) Run(fl *flag.FlagSet) {
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