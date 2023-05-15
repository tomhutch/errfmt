package funcfromgetter

import (
	"fmt"
	"os"
)

type fnGetter struct {
	GetFn func(string) error
}

func okFormat() error {
	f := fnGetter{GetFn: os.Chdir}
	err := f.GetFn("")
	if err != nil {
		return fmt.Errorf("f.GetFn: %w", err)
	}

	return nil
}

func badFormat() error {
	f := fnGetter{GetFn: os.Chdir}
	err := f.GetFn("")
	if err != nil {
		return fmt.Errorf("failed to GetFn: %w", err) // want `error message not prefixed in expected format`
	}

	return nil
}
