package go_validator

import (
"context"
"fmt"
"strings"
)

// IntValidator is a validator that checks if the value is an integer of any type.
func IntValidator(_ context.Context, value any) (twigBlock bool, errs []error) {

	switch value.(type) {
	case float64:
		s := fmt.Sprintf("%f", value)
		dec := strings.Split(s, ".")[1]
		if strings.Count(dec, "0") == len(dec) {
			return
		}
	case int:
		return
	}

	return true, []error{NotAnIntegerError{}}
}
