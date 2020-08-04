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
	// Mock for root/parent command.
	mockParentCmd struct{}
	// Mock subcommand with aliases and a nested sucommand of its own.
	mockSubCmd struct {
		buf *bytes.Buffer
	}
	// Mock subcommand with aliases and no nested subcommands.
	mockSubCmdNoNested struct{}
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

func (c *mockSubCmd) Subcommands() []Command {
	return []Command{new(mockSubCmd)}
}

func (c *mockSubCmd) WriteString(s string) (int, error) {
	return c.buf.WriteString(s)
}

func (c *mockSubCmd) Spec() CommandSpec {
	return CommandSpec{
		Name:    "mockSubCmd",
		Usage:   "Test a subcommand.",
		Aliases: []string{"s", "sc", "sub"},
		Desc:    "A simple mock subcommand with aliases and its own subcommand.",
	}
}

func (c *mockSubCmdNoNested) Run(fl *pflag.FlagSet) {}

func (c *mockSubCmdNoNested) Spec() CommandSpec {
	return CommandSpec{
		Name:    "mockSubCmdNoNested",
		Usage:   "Used for help output tests.",
		Aliases: []string{"s", "sc", "sub"},
		Desc:    "A simple mock subcommand with aliases and no nested subcommands.",
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

Description: Mock parent command description.

Commands:
	s, sc, sub, mockSubCmd  - A simple mock subcommand with aliases and its own subcommand.
`
		buf := new(bytes.Buffer)
		cmd := new(mockParentCmd)
		name := cmd.Spec().Name
		fl := pflag.NewFlagSet(name, pflag.ExitOnError)
		// If the help output doesn't contain the subcommand and
		// isn't formatted the way we expect the test will fail.
		renderHelp(buf, name, cmd, fl)
		got := buf.String()
		assert.Equal(t, t.Name(), expected, got)
	})
}

func TestSubCmdHelpOutput(t *testing.T) {
	withNested := `Usage: mockSubCmd Test a subcommand.

Aliases: s, sc, sub

Description: A simple mock subcommand with aliases and its own subcommand.

Commands:
	s, sc, sub, mockSubCmd  - A simple mock subcommand with aliases and its own subcommand.
`

	noNested := `Usage: mockSubCmdNoNested Used for help output tests.

Aliases: s, sc, sub

Description: A simple mock subcommand with aliases and no nested subcommands.
`

	for _, test := range []struct {
		name, expected string
		cmd            Command
	}{
		{
			name:     "subcmd w/nested subcmd.",
			expected: withNested,
			cmd:      new(mockSubCmd),
		},
		{
			name:     "subcmd w/no nested subcmds.",
			expected: noNested,
			cmd:      new(mockSubCmdNoNested),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			name := test.cmd.Spec().Name
			fl := pflag.NewFlagSet(name, pflag.ExitOnError)
			// If the help output is not written to the buffer
			// in the format we expect then the test will fail.
			renderHelp(buf, name, test.cmd, fl)
			got := buf.String()
			assert.Equal(t, t.Name(), test.expected, got)
		})
	}
}
