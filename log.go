package main

import (
	"fmt"
	"os"
)

var verbose bool

func debugf(format string, args ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%s", format), args...)
	}
}

func debug(args ...interface{}) {
	if verbose {
		fmt.Fprintln(os.Stderr, args...)
	}
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("ERROR: %s", format), args...)
	os.Exit(1)
}

func fatal(args ...interface{}) {
	fmt.Fprint(os.Stderr, "ERROR: ")
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
}
