package utils

import "go.uber.org/multierr"

func ReturnThunk[T any](f func() ([]T, []error)) ([]T, error) {
	ret, errs := f()
	if errs != nil {
		return nil, multierr.Combine(errs...)
	}
	return ret, nil
}
