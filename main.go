package main

import (
	"flag"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	// Don't use it: just to not crash on -unsafeptr flag from go vet
	flag.Bool("unsafeptr", false, "")

	cfg := NewDefaultConfig()
	singlechecker.Main(NewAnalyzer(cfg))
}
