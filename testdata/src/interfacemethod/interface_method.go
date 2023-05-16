package interfacemethod

import (
	"encoding/json"
	"fmt"
	"strings"
)

type errorer interface {
	Decode(v interface{}) error
}

type foo struct {
	bar errorer
}

func okFormat() error {
	d := json.NewDecoder(strings.NewReader("hello world"))
	f := foo{d}
	var str string
	err := f.bar.Decode(&str)
	if err != nil {
		return fmt.Errorf("f.bar.Decode: %w", err)
	}

	return nil
}

func badFormat() error {
	d := json.NewDecoder(strings.NewReader("hello world"))
	f := foo{d}
	var str string
	err := f.bar.Decode(&str)
	if err != nil {
		return fmt.Errorf("failed to decode: %w", err) // want `error message not prefixed in expected format`
	}

	return nil
}
