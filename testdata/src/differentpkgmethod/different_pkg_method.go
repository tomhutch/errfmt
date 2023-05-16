package differentpkgmethod

import (
	"fmt"
	"io"
)

func okFormat() error {
	writer := io.NewOffsetWriter(nil, 0)
	_, err := writer.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("writer.Seek: %w", err)
	}

	return nil
}

func badFormat() error {
	writer := io.NewOffsetWriter(nil, 0)
	_, err := writer.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("io.writer.Seek: %w", err) // want `error message not prefixed in expected format`
	}

	return nil
}
