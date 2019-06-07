package cli

import (
	"flag"
	"fmt"
	"io"
	"text/tabwriter"
	"unicode/utf8"
)

func flagDashes(name string) string {
	if utf8.RuneCountInString(name) > 1 {
		return "--"
	}
	return "-"
}

func renderFlagHelp(fl *flag.FlagSet, w io.Writer) {
	var count int
	fl.VisitAll(func(f *flag.Flag) {
		if count == 0 {
			fmt.Fprintf(w, "\n%v flags:\n", fl.Name())
		}

		count++
		if f.DefValue == "" {
			fmt.Fprintf(w, "\t%v%v\t%v\n", flagDashes(f.Name), f.Name, f.Usage)
		} else {
			fmt.Fprintf(w, "\t%v%v\t%v\t(%v)\n", flagDashes(f.Name), f.Name, f.Usage, f.DefValue)
		}
	})
}

// renderHelp generates a command's help page.
func renderHelp(cmd Command, fl *flag.FlagSet, w io.Writer) {
	// Render usage and description.
	fmt.Fprintf(w, "Usage: %v %v\n\n",
		fl.Name(), cmd.Spec().Usage,
	)
	fmt.Fprintf(w, "%v\n", cmd.Spec().Desc)

	// Render flag help.
	renderFlagHelp(fl, w)

	// Render subcommand summaries.
	pc, ok := cmd.(ParentCommand)
	if ok {
		if len(pc.Subcommands()) > 0 {
			// Give some space from flags.
			fmt.Fprintf(w, "\n")
			fmt.Fprint(w, "Commands:\n")
		}

		tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', tabwriter.StripEscape)
		for _, cmd := range pc.Subcommands() {
			if cmd.Spec().Hidden {
				continue
			}

			// \xFF is used to escape the leading \t so tabwriter ignores it
			fmt.Fprintf(tw,
				"\xFF\t\xFF%v\t%v\n",
				cmd.Spec().Name,
				cmd.Spec().ShortDesc(),
			)
		}

		tw.Flush()
	}
}
