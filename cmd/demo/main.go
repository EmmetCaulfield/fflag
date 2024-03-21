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
	fflag.Var(&help, '?', "help", "print a help message")
	fflag.Var(&nocase, 'i', "ignore-case", "ignore case in patterns",
		fflag.WithAlias('y', "", true), fflag.NotImplemented())
	fflag.Var(&foo, 0, "foo", "test setting a string", fflag.WithRepeats(true))
	fflag.Var(&bar, 'b', "bar", "test setting an int slice")
	afds := fflag.CommandLine.AlignedFlagDescriptions("  ", "  ", "")
	fmt.Println(strings.Join(afds, "\n"))
	fflag.Parse()
	fmt.Println()
	fflag.CommandLine.DumpFlags()
	return
}
