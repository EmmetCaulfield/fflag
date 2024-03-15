package main

import (
	"fmt"
	"os"

	"github.com/EmmetCaulfield/fflag"
)

func main() {
	fmt.Printf("%d\n----\n", len(os.Args))
	var help bool
	fflag.Var(&help, "help", "print a help message", fflag.WithShortcut('?'), fflag.WithTypeTag("BOOL"))
	fmt.Println(fflag.Lookup("help").FlagString())
	fflag.Parse()
	return
}
