package utils

import (
	"errors"
)

func AreErrorsEqual(actual, expected error) bool {
	if actual != nil && expected != nil {
		return actual.Error() == expected.Error()
	}
	return errors.Is(actual, expected)
}
