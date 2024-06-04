package validator

import (
	"fmt"
	"regexp"
)

func ValidateName(arg string) error {
	if regexp.MustCompile(`^[a-zA-Z0-9_\-]{1,50}$`).MatchString(arg) {
		return nil
	}
	return fmt.Errorf("`%s` must comprise letters, numbers, underscore, dash and not have more than 50 characters", arg)
}
