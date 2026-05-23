// Package picker wraps go-fuzzyfinder behind a simple interface so that
// commands can be tested with a mock implementation.
package picker

import (
	"fmt"

	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
)

// Picker selects one item from a list.
type Picker interface {
	// Pick presents items to the user and returns the selected item.
	// Returns ("", ErrAborted) if the user cancels.
	Pick(prompt string, items []string) (string, error)
}

// ErrAborted is returned when the user cancels the picker (e.g. presses Esc).
var ErrAborted = fmt.Errorf("selection aborted")

// FuzzyPicker is the real implementation backed by go-fuzzyfinder.
type FuzzyPicker struct{}

// New returns a new FuzzyPicker.
func New() Picker {
	return &FuzzyPicker{}
}

// Pick opens the interactive fuzzy finder and returns the chosen item.
func (p *FuzzyPicker) Pick(prompt string, items []string) (string, error) {
	if len(items) == 0 {
		return "", fmt.Errorf("no items to pick from")
	}

	idx, err := fuzzyfinder.Find(
		items,
		func(i int) string { return items[i] },
		fuzzyfinder.WithPromptString(prompt+" "),
	)
	if err != nil {
		if err == fuzzyfinder.ErrAbort {
			return "", ErrAborted
		}
		return "", err
	}

	return items[idx], nil
}
