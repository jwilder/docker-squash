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
	if verbose {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("ERROR: %s", format), args...)
		signals <- os.Interrupt
		wg.Wait()
		os.Exit(1)
	}
}

func fatal(args ...interface{}) {
	if verbose {
		fmt.Fprint(os.Stderr, "ERROR: ")
		fmt.Fprintln(os.Stderr, args...)
		signals <- os.Interrupt
		wg.Wait()
		os.Exit(1)
	}
}
