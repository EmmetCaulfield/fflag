package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/EmmetCaulfield/fflag"
)

func main() {
	fmt.Printf("%d\n----\n", len(os.Args))
	var help, nocase bool
	var foo string
	bar := []int{}
	fflag.Var(&help, "help", "print a help message", fflag.WithShortcut('?'))
	fflag.Var(&nocase, "ignore-case", "ignore case in patterns",
		fflag.WithShortcut('i'), fflag.WithAlias("", 'y', true), fflag.NotImplemented())
	fflag.Var(&foo, "foo", "test setting a string")
	fflag.Var(&bar, "bar", "test setting an int slice", fflag.WithShortcut('b'))
	afds := fflag.CommandLine.AlignedFlagDescriptions("  ", "  ", "")
	fmt.Println(strings.Join(afds, "\n"))
	fflag.Parse()
	fmt.Println()
	fflag.CommandLine.DumpFlags()
	return
}
