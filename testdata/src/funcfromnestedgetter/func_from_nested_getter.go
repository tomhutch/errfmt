package funcfromnestedgetter

import (
	"fmt"
	"os"
)

type errFn func(string) error

type fnGetter struct {
	GetFn func() errFn
}

func okFormat() error {
	nestedGetter := func(s string) error { return os.Chdir(s) }
	f := fnGetter{GetFn: func() errFn { return nestedGetter }}
	err := f.GetFn()("")
	if err != nil {
		return fmt.Errorf("f.GetFn: %w", err)
	}

	return nil
}

func badFormat() error {
	nestedGetter := func(s string) error { return os.Chdir(s) }
	f := fnGetter{GetFn: func() errFn { return nestedGetter }}
	err := f.GetFn()("")
	if err != nil {
		return fmt.Errorf("failed to GetFn: %w", err) // want `error message not prefixed in expected format`
	}

	return nil
}
