package config

import "fmt"

type ValueNotFoundError struct {
	key string
}

func (err *ValueNotFoundError) Error() string {
	return fmt.Sprintf("[cound not find config value by key: %s]", err.key)
}

type ValueTypeNotMatchError struct {
	key      string
	value    any
	expected string
}

func (err *ValueTypeNotMatchError) Error() string {
	return fmt.Sprintf(
		"[unexpected config value type] key: %s; value: %v(%T); expected: %s",
		err.key,
		err.value,
		err.value,
		err.expected,
	)
}

type InvalidTargetTypeError struct {
	t string
}

func (err *InvalidTargetTypeError) Error() string {
	return fmt.Sprintf(
		"invalid target type, expected ptr, got %s", err.t,
	)
}
