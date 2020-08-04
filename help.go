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
func renderHelp(w io.Writer, fullName string, cmd Command, fl *pflag.FlagSet) {
	var b strings.Builder
	spec := cmd.Spec()
	fmt.Fprintf(&b, "Usage: %v %v\n\n", fullName, spec.Usage)

	// If the command has aliases, add them to the output as a comma-separated list.
	if spec.HasAliases() {
		fmt.Fprintf(&b, "Aliases: %s\n\n", strings.Join(spec.Aliases, ", "))
	}
	// Print usage and description.
	fmt.Fprintf(w, "%sDescription: %s\n", b.String(), spec.Desc)

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
			spec := cmd.Spec()

			if spec.Hidden {
				continue
			}

			allNames := strings.Join(append(spec.Aliases, spec.Name), ", ")

			fmt.Fprintf(tw,
				tabEscape+"\t"+tabEscape+"%v\t- %v\n",
				allNames,
				spec.ShortDesc(),
			)
		}
	}
}
