package cli

import (
	"bytes"
	"os"
	"testing"

	"cdr.dev/slog/sloggers/slogtest/assert"
	"github.com/spf13/pflag"
)

var subCmd = new(mockSubCmd)

type (
	mockParentCmd struct{}
	mockSubCmd    struct {
		buf *bytes.Buffer
	}
)

func (c *mockParentCmd) Run(fl *pflag.FlagSet) {}

func (c *mockParentCmd) Subcommands() []Command {
	return []Command{subCmd}
}

func (c *mockParentCmd) Spec() CommandSpec {
	return CommandSpec{
		Name:  "mockParentCmd",
		Usage: "Mock parent command usage.",
		Desc:  "Mock parent command description.",
	}
}

func (c *mockSubCmd) Run(fl *pflag.FlagSet) {
	c.buf = new(bytes.Buffer)
	_, err := c.WriteString("success")
	if err != nil {
		fl.Usage()
	}
}

func (c *mockSubCmd) WriteString(s string) (int, error) {
	return c.buf.WriteString(s)
}

func (c *mockSubCmd) Spec() CommandSpec {
	return CommandSpec{
		Name:    "mockSubCmd",
		Usage:   "Test a subcommand.",
		Aliases: []string{"s", "sc", "sub"},
		Desc:    "A simple mock subcommand with aliases.",
	}
}

func TestSubCmdAliases(t *testing.T) {
	for _, alias := range []string{"s", "sc", "sub"} {
		t.Run(alias, func(t *testing.T) {
			// Setup command.
			cmd := new(mockParentCmd)
			os.Args = []string{cmd.Spec().Name, alias}
			// Run command.
			RunRoot(cmd)
			// If "success" isn't written into the buffer
			// then we failed to find the subcommand by alias.
			got := string(subCmd.buf.Bytes())
			assert.Equal(t, t.Name(), "success", got)
		})
	}
}

func TestCmdHelpOutput(t *testing.T) {
	t.Run(t.Name(), func(t *testing.T) {
		expected := `Usage: mockParentCmd Mock parent command usage.

Mock parent command description.

Commands:
	s,sc,sub,mockSubCmd  - A simple mock subcommand with aliases.
`
		buf := new(bytes.Buffer)
		cmd := new(mockParentCmd)
		name := cmd.Spec().Name
		fl := pflag.NewFlagSet(name, pflag.ExitOnError)
		// If the help output doesn't contain the subcommand and
		// isn't formatted the way we expect the test will fail.
		renderHelp(name, cmd, fl, buf)
		got := string(buf.Bytes())
		assert.Equal(t, t.Name(), expected, got)
	})
}
