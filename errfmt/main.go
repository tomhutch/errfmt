package main

import (
	"flag"

	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/tomhutch/errfmt/core"
)

func main() {
	// Don't use it: just to not crash on -unsafeptr flag from go vet
	flag.Bool("unsafeptr", false, "")

	cfg := core.NewDefaultConfig()
	singlechecker.Main(core.NewAnalyzer(cfg))
}
