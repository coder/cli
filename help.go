package cli

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"unicode"
	"unicode/utf8"

	"github.com/spf13/pflag"
)

func flagDashes(name string) string {
	if utf8.RuneCountInString(name) > 1 {
		return "--"
	}
	return "-"
}

// fmtDefValue adds quotes around default value strings that contain spaces so
// the help representation matches what you would need to do when running a
// command.
func fmtDefValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		if strings.IndexFunc(v, unicode.IsSpace) != -1 {
			return fmt.Sprintf(`"%v"`, v)
		}
		return v

	default:
		return fmt.Sprintf("%v", v)
	}
}

// \xFF is used to escape a \t so tabwriter ignores it.
const tabEscape = "\xFF"

func renderFlagHelp(fullName string, fl *pflag.FlagSet, w io.Writer) {
	if fl.HasFlags() {
		fmt.Fprintf(w, "\n%s flags:\n", fullName)
		fmt.Fprint(w, fl.FlagUsages())
	}
}

// renderHelp generates a command's help page.
func renderHelp(fullName string, cmd Command, fl *pflag.FlagSet, w io.Writer) {
	// Render usage and description.
	fmt.Fprintf(w, "Usage: %v %v\n\n",
		fullName, cmd.Spec().Usage,
	)
	fmt.Fprintf(w, "%v\n", cmd.Spec().Desc)

	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', tabwriter.StripEscape)
	defer tw.Flush()

	// Render flag help.
	renderFlagHelp(fullName, fl, tw)

	// Render subcommand summaries.
	pc, ok := cmd.(ParentCommand)
	if ok {
		if len(pc.Subcommands()) > 0 {
			// Give some space from flags.
			fmt.Fprintf(w, "\n")
			fmt.Fprint(w, "Commands:\n")
		}

		for _, cmd := range pc.Subcommands() {
			if cmd.Spec().Hidden {
				continue
			}

			fmt.Fprintf(tw,
				tabEscape+"\t"+tabEscape+"%v\t%v\n",
				cmd.Spec().Name,
				cmd.Spec().ShortDesc(),
			)
		}
	}
}
