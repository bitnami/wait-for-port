package main

import (
	"flag"
	"os"
	"regexp"
	"testing"

	cassert "github.com/juamedgod/cliassert"
)

func TestHelp(t *testing.T) {
	r := RunTool("--help")
	r.AssertErrorMatch(t, regexp.MustCompile(`(?ms).*Usage.*--host.*`))
}

func TestMain(m *testing.M) {
	if os.Getenv("BE_TOOL") == "1" {
		main()
		os.Exit(0)
		return
	}
	flag.Parse()
	c := m.Run()
	os.Exit(c)
}

func RunTool(args ...string) cassert.CmdResult {
	os.Setenv("BE_TOOL", "1")
	defer os.Unsetenv("BE_TOOL")
	return cassert.ExecCommand(os.Args[0], args...)
}
