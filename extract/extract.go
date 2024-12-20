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

// ErrIsNotStruct if the object given is not a struct.
var ErrIsNotStruct = errors.New("object given is not struct")
