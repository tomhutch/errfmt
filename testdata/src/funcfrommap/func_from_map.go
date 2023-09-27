package funcfrommap

import (
	"fmt"
)

type stadium struct{}

func (s stadium) Play() error {
	return nil
}

func okFormat() error {
	allStadiums := map[string]stadium{
		"Cambridge": {},
	}
	err := allStadiums["Cambridge"].Play()
	if err != nil {
		return fmt.Errorf("allStadiums[Cambridge].Play: %w", err)
	}

	return nil
}

func okFormatNested() error {
	allVenues := map[string]map[string]stadium{
		"stadiums": {
			"Cambridge": {},
		},
	}
	err := allVenues["stadiums"]["Cambridge"].Play()
	if err != nil {
		return fmt.Errorf("allVenues[stadiums][Cambridge].Play: %w", err)
	}

	return nil
}

func badFormat() error {
	allStadiums := map[string]stadium{
		"Cambridge": {},
	}
	err := allStadiums["Cambridge"].Play()
	if err != nil {
		return fmt.Errorf("Play: %w", err) // want `error message not prefixed in expected format`
	}

	return nil
}
