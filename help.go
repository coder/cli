package cli

import (
	"flag"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"unicode"
	"unicode/utf8"
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

func renderFlagHelp(fl *flag.FlagSet, w io.Writer) {
	var count int
	fl.VisitAll(func(f *flag.Flag) {
		if count == 0 {
			fmt.Fprintf(w, "\n%v flags:\n", fl.Name())
		}

		count++
		if f.DefValue == "" {
			fmt.Fprintf(w, tabEscape+"\t%v%v\t%v\n", flagDashes(f.Name), f.Name, f.Usage)
		} else {
			fmt.Fprintf(w, tabEscape+"\t%v%v=%v\t%v\n", flagDashes(f.Name), f.Name, fmtDefValue(f.DefValue), f.Usage)
		}
	})
}

// renderHelp generates a command's help page.
func renderHelp(s *Config,  w io.Writer) {
	// Render usage and description.
	fmt.Fprintf(w, "Usage: %v %v\n\n",
		s.Flags.Name(), s.Usage,
	)
	fmt.Fprintf(w, "%v\n", s.Desc)

	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', tabwriter.StripEscape)
	defer tw.Flush()

	// Render flag help.
	renderFlagHelp(s.Flags, tw)

	// Render subcommand summaries.
	if len(s.Subcommands) > 0 {
		// Give some space from flags.
		fmt.Fprintf(w, "\n")
		fmt.Fprint(w, "Commands:\n")
	}

	for _, cmd := range s.Subcommands {
		sc := commandConfig(cmd)
		if sc.Hidden {
			continue
		}

		fmt.Fprintf(tw,
			tabEscape+"\t"+tabEscape+"%v\t%v\n",
			sc.Name,
			sc.ShortDesc(),
		)
	}
}
