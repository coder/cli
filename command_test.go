package cli

import (
	"bytes"
	"os"
	"testing"

	"cdr.dev/slog/sloggers/slogtest/assert"
	"github.com/spf13/pflag"
)

const (
	success = "test successful"

	expectedParentCmdHelpOutput = `Usage: mockParentCmd Mock parent command usage.

Mock parent command description.

Commands:
	s,sc,sub,mockSubCmd  - A simple mock subcommand with aliases.
`
)

var subCmd = &mockSubCmd{buf: new(bytes.Buffer)}

type (
	mockParentCmd struct{}
	mockSubCmd    struct{ buf *bytes.Buffer }
)

func (c *mockParentCmd) Run(fl *pflag.FlagSet) {}

func (c *mockParentCmd) Subcommands() []Command { return []Command{subCmd} }

func (c *mockParentCmd) Spec() CommandSpec {
	return CommandSpec{
		Name:  "mockParentCmd",
		Usage: "Mock parent command usage.",
		Desc:  "Mock parent command description.",
	}
}

func (c *mockSubCmd) Run(fl *pflag.FlagSet) {
	c.buf = new(bytes.Buffer)
	if _, err := c.Write([]byte(success)); err != nil {
		fl.Usage()
	}
}

func (c *mockSubCmd) Write(b []byte) (int, error) { return c.buf.Write(b) }

func (c *mockSubCmd) Spec() CommandSpec {
	return CommandSpec{
		Name:    "mockSubCmd",
		Usage:   "Test a subcommand.",
		Aliases: []string{"s", "sc", "sub"},
		Desc:    "A simple mock subcommand with aliases.",
	}
}

func TestSubCmdAliases(t *testing.T) {
	for _, test := range []struct {
		name, expected string
	}{
		{
			name:     "s",
			expected: success,
		},
		{
			name:     "sc",
			expected: success,
		},
		{
			name:     "sub",
			expected: success,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			// Since the alias is the name of the test
			// we can just pass it as the alias arg.
			os.Args = []string{"mockParentCmd", test.name}
			// Based on os.Args, when splitArgs is invoked,
			// it should be able to deduce the subcommand we want
			// based on the new alias map it's being passed.
			RunRoot(&mockParentCmd{})
			// The success const is never written into the buffer
			// if the subcommand fails to be invoked by alias.
			got := string(subCmd.buf.Bytes())
			assert.Equal(t, test.name, test.expected, got)
		})
	}
}

func TestSubcmdAliasesInParentCmdHelpOutput(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := &mockParentCmd{}
	name := cmd.Spec().Name
	fl := pflag.NewFlagSet(name, pflag.ExitOnError)
	// If the help output is not written to the buffer
	// in the format we expect then the test will fail.
	renderHelp(name, cmd, fl, buf)
	got := string(buf.Bytes())
	expected := expectedParentCmdHelpOutput
	assert.Equal(t, "display_subcmd_aliases_in_parentcmd_help_output", expected, got)
}
