// Package cli provides a thin CLI abstraction around the standard flag package.
// It is minimal, command-struct-oriented, and trades off "power" for flexibility
// and clarity at the caller level.
//
// It pretends that Go's single dash (-flag) support doesn't exist, and renders
// helps with --.
//
// Optional interface adherence can be asserted with a statement like
//  var _ interface {
//  	cli.Command
//  	cli.FlaggedCommand
//  	cli.ParentCommand
//  } = new(rootCmd)
//
// This package is exported since we think it could be interesting to the community,
// but the API is subject to change in support of sail.
package cli
