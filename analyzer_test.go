package main

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAll(t *testing.T) {
	tests := []struct {
		name        string
		packageName string
	}{
		{name: "different package func", packageName: "differentpkgfunc"},
		{name: "different package method", packageName: "differentpkgmethod"},
		{name: "error not wrapped", packageName: "errnotwrapped"},
		{name: "func from getter", packageName: "funcfromgetter"},
		{name: "func from nested getter", packageName: "funcfromnestedgetter"},
		{name: "func literal", packageName: "funcliteral"},
		{name: "interface method", packageName: "interfacemethod"},
		{name: "recievermethod method", packageName: "recievermethod"},
		{name: "var assigned func", packageName: "varassignedfunc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := NewDefaultConfig()

			analysistest.Run(
				t,
				analysistest.TestData(),
				NewAnalyzer(cfg),
				tt.packageName,
			)
		})
	}
}
