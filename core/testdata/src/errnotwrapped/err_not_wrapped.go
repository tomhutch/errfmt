package errnotwrapped

import (
	"fmt"
)

func skippedFormat() error {
	_, err := fmt.Scanf("failed to do something")
	if err != nil {
		return err
	}

	return nil
}
