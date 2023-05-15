package varassignedfunc

import (
	"fmt"
)

func okFormat() error {
	funcName := func() error { return nil }
	err := funcName()
	if err != nil {
		return fmt.Errorf("funcName: %w", err)
	}

	return nil
}

func badFormat() error {
	funcName := func() error { return nil }
	err := funcName()
	if err != nil {
		return fmt.Errorf("failure in funcName: %w", err) // want `error message not prefixed in expected format`
	}

	return nil
}
