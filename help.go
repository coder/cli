package cli

import (
	"flag"
	"fmt"
	"io"
	"unicode/utf8"
)

func flagDashes(name string) string {
	if utf8.RuneCountInString(name) > 1 {
		return "--"
	}
	return "-"
}

func renderFlagHelp(fl *flag.FlagSet, w io.Writer) {
	fmt.Fprintf(w, "%v flags:\n", fl.Name())
	var count int
	fl.VisitAll(func(f *flag.Flag) {
		count++
		if f.DefValue == "" {
			fmt.Fprintf(w, "\t%v%v\t%v\n", flagDashes(f.Name), f.Name, f.Usage)
		} else {
			fmt.Fprintf(w, "\t%v%v\t%v\t(%v)\n", flagDashes(f.Name), f.Name, f.Usage, f.DefValue)
		}
	})
	if count == 0 {
		fmt.Fprintf(w, "\n")
	}
}

// renderHelp generates a command's help page.
func renderHelp(cmd Command, fl *flag.FlagSet, w io.Writer) {
	// Render usage and description.
	fmt.Fprintf(w, "Usage: %v %v\n\n",
		fl.Name(), cmd.Spec().Usage,
	)
	fmt.Fprintf(w, "%v\n\n", cmd.Spec().Desc)

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
		for _, cmd := range pc.Subcommands() {
			fmt.Fprintf(w, "\t%v\t%v\n", cmd.Spec().Name, cmd.Spec().ShortDesc())
		}
	}
}
