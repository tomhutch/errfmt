package testdata_original

import "fmt"

func firstFn() error {
	_, err := fmt.Scanf("1")
	if err != nil {
		return fmt.Errorf("failed to scanf: %w", err)
	}

	return nil
}

func secondFn() error {
	_, err := fmt.Scanf("2")
	if err != nil {
		return fmt.Errorf("scanf: %w", err)
	}

	return nil
}

func thirdFn() error {
	_, err := fmt.Scanf("3")
	if err != nil {
		return fmt.Errorf("failed to scanf: %w", err)
	}

	return nil
}
