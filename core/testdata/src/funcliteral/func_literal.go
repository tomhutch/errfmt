package funcliteral

import "fmt"

func skippedFormat() error {
	err := func() error {
		return nil
	}()
	if err != nil {
		return fmt.Errorf("function literal failed: %w", err)
	}

	return nil
}
