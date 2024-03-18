package main

import (
	"fmt"
	"os"

	"github.com/EmmetCaulfield/fflag"
)

func main() {
	fmt.Printf("%d\n----\n", len(os.Args))
	var help, nocase bool
	fflag.Var(&help, "help", "print a help message", fflag.WithShortcut('?'))
	fflag.Var(&nocase, "ignore-case", "ignore case", fflag.WithShortcut('i'), fflag.WithAlias("", 'y', true))
	for _, f := range fflag.CommandLine.FlagList {
		fmt.Println(f.FlagString(), "\t", f.DescString())
	}
	fflag.Parse()
	return
}
