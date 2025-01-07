package extract

import "errors"

type (
	queryType        string
	queryPlaceholder string
)

var (
	// CreateType create queryType value.
	CreateType queryType = "create"
	// UpdateType update queryType value.
	UpdateType queryType = "update"

	// QuestionMarkPlaceholder question mark queryPlaceholder value.
	QuestionMarkPlaceholder queryPlaceholder = "?"
	// DollarPlaceholder dollar queryPlaceholder value.
	DollarPlaceholder queryPlaceholder = "$"
)

var (
	// ErrIsNotStruct if the object given is not a struct.
	ErrIsNotStruct = errors.New("object given is not struct")

	// ErrBaseQuery if the final query is empty.
	ErrBaseQuery = errors.New("final query is same the base")
)
