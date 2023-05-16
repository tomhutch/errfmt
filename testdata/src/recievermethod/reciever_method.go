package recievermethod

import (
	"errors"
	"fmt"
)

type errorer interface {
	Decode(v interface{}) error
}

type foo struct{}

var _ errorer = (*foo)(nil)

func (f *foo) Decode(v interface{}) error {
	return errors.New("something broke")
}

func okFormat() error {
	f := foo{}
	var str string
	err := f.Decode(&str)
	if err != nil {
		return fmt.Errorf("f.Decode: %w", err)
	}

	return nil
}

func badFormat() error {
	f := foo{}
	var str string
	err := f.Decode(&str)
	if err != nil {
		return fmt.Errorf("failed to decode: %w", err) // want `error message not prefixed in expected format`
	}

	return nil
}
