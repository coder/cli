// Package cli provides a thin CLI abstraction around the standard flag package.
// It is minimal, command-struct-oriented, and trades off "power" for flexibility
// and clarity at the caller level.
//
// It pretends that Go's single dash (-flag) support doesn't exist, and renders
// helps with --.
package cli
