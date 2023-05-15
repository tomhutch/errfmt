package main

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/tomhutch/errfmt/core"
)

func TestFixer(t *testing.T) {
	dir := analysistest.TestData()
	sourceFilePath := dir + "_original/src.go"
	sourceFilePath2 := dir + "_original/src2.go"
	copyFilePath := dir + "/out.go"
	copyFilePath2 := dir + "/out2.go"
	err := copyFile(sourceFilePath, copyFilePath)
	if err != nil {
		t.Fatal("copyFile failed")
	}
	err = copyFile(sourceFilePath2, copyFilePath2)
	if err != nil {
		t.Fatal("copyFile failed")
	}

	cfg := core.NewDefaultConfig()
	cfg.EnableFixer = true

	analysistest.Run(
		t,
		analysistest.TestData(),
		core.NewAnalyzer(cfg),
		"",
	)

	data, err := os.ReadFile(copyFilePath)
	if err != nil {
		t.Fatal("read failed")
	}
	expectedContents := expected()
	if expectedContents != string(data) {
		t.Log(cmp.Diff(expectedContents, string(data)))
		t.Fail()
	}

	data2, err := os.ReadFile(copyFilePath2)
	if err != nil {
		t.Fatal("read failed")
	}
	expectedContents2 := expected2()
	if expectedContents2 != string(data2) {
		t.Log(cmp.Diff(expectedContents2, string(data2)))
		t.Fail()
	}
}

func copyFile(src string, dst string) error {
	// Read all content of src to data, may cause OOM for a large file.
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	// Write data to dst
	err = os.WriteFile(dst, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func expected() string {
	return `package testdata_original

import "fmt"

func firstFn() error {
	_, err := fmt.Scanf("1")
	if err != nil {
		return fmt.Errorf("fmt.Scanf: %w", err)
	}

	return nil
}

func secondFn() error {
	_, err := fmt.Scanf("2")
	if err != nil {
		return fmt.Errorf("fmt.Scanf: %w", err)
	}

	return nil
}

func thirdFn() error {
	_, err := fmt.Scanf("3")
	if err != nil {
		return fmt.Errorf("fmt.Scanf: %w", err)
	}

	return nil
}
`
}

func expected2() string {
	return `package testdata_original

import "fmt"

func fourthFn() error {
	_, err := fmt.Scanf("4")
	if err != nil {
		return fmt.Errorf("fmt.Scanf: %w", err)
	}

	return nil
}

func fifthFn() error {
	_, err := fmt.Scanf("5")
	if err != nil {
		return fmt.Errorf("fmt.Scanf: %w", err)
	}

	return nil
}
`
}
