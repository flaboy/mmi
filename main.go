package main

import (
	"flag"
	"fmt"
	"github.com/flaboy/mmi/parser"
	"os"
)

var (
	index_mode   bool
	server_mode  bool
	help_mode    bool
	summary_mode bool
	latex_mode   bool
	depth        int64
	workdir      string
)

func main() {

	flag.BoolVar(&summary_mode, "summary", false, "build SUMMARY.md")
	flag.BoolVar(&index_mode, "json", false, "build index.json")
	flag.BoolVar(&latex_mode, "latex", false, "output latex")
	flag.BoolVar(&server_mode, "server", false, "start server")
	flag.BoolVar(&help_mode, "help", false, "show help")
	flag.Int64Var(&depth, "depth", 5, "depth")
	flag.Parse()

	args := flag.Args()

	workdir = "."
	if len(args) > 0 {
		workdir = args[0]
	}

	if server_mode {
		start_server(workdir)
		return
	}

	n := parser.Open(workdir)
	if index_mode {
		n.UpdateJson()
		return
	}

	if latex_mode {
		n.ToLatex()
		return
	}

	if summary_mode {
		if depth < 1 {
			depth = 5
		}
		n.UpdateRummary(depth)
		return
	}

	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}
