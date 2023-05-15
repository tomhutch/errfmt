package differentpkgfunc

import (
	"fmt"
)

func okFormat() error {
	_, err := fmt.Scanf("")
	if err != nil {
		return fmt.Errorf("fmt.Scanf: %w", err)
	}

	return nil
}

func badFormat() error {
	_, err := fmt.Scanf("")
	if err != nil {
		return fmt.Errorf("fmt.Scanf: %w", err) // want `error message not prefixed in expected format`
	}

	return nil
}
