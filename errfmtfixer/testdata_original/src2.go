package testdata_original

import "fmt"

func fourthFn() error {
	_, err := fmt.Scanf("4")
	if err != nil {
		return fmt.Errorf("scanf: %w", err)
	}

	return nil
}

func fifthFn() error {
	_, err := fmt.Scanf("5")
	if err != nil {
		return fmt.Errorf("failed to scanf: %w", err)
	}

	return nil
}
