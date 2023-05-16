package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	cfg := NewDefaultConfig()
	singlechecker.Main(NewAnalyzer(cfg))
}
