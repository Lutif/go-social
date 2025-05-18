package customerror

import "errors"

func CheckForCustomErr(err, target, replacer error) error {
	switch {
	case errors.Is(err, target):
		return replacer
	default:
		return err
	}
}
