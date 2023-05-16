package funcfromnestedgetter

import (
	"fmt"
	"os"
)

type fnGetter struct {
	GetFn func() func(string) error
}

func okFormat() error {
	nestedGetter := func() func(string) error { return os.Chdir }
	f := fnGetter{GetFn: nestedGetter}
	err := f.GetFn()("")
	if err != nil {
		return fmt.Errorf("f.GetFn: %w", err)
	}

	return nil
}

func badFormat() error {
	nestedGetter := func() func(string) error { return os.Chdir }
	f := fnGetter{GetFn: nestedGetter}
	err := f.GetFn()("")
	if err != nil {
		return fmt.Errorf("failed to GetFn: %w", err) // want `error message not prefixed in expected format`
	}

	return nil
}
