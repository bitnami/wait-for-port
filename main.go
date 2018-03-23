package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
)

func main() {
	cmd := NewWaitForPortCmd()
	parser := flags.NewParser(cmd, flags.HelpFlag|flags.PassDoubleDash)

	args, err := parser.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := cmd.Execute(args); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
