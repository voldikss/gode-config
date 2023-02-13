package parser

import "fmt"

type InvalidValueError struct {
	Key  string
	Kind string
}

type KeyNotFoundError struct {
	key string
}

type ParseError struct {
	Parser string
}

func (e *InvalidValueError) Error() string {
	return fmt.Sprintf("value for the key %q is not a %s", e.Key, e.Kind)
}

func (e *KeyNotFoundError) Error() string {
	return fmt.Sprintf("not found by key %s", e.key)
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse failed using parser %s", e.Parser)
}
