package main

import (
	"flag"
	"github.com/flaboy/mmi/parser"
)

func main() {
	var toc_mode bool
	flag.BoolVar(&toc_mode, "json", false, "json mode")
	flag.Parse()

	args := flag.Args()

	if len(args) > 0 {
		n := parser.Open(args[0])
		if toc_mode {
			n.UpdateJson()
		} else {
			n.UpdateReadme(2)
		}
	}
}
