package extract

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type (
	queryType        string
	queryPlaceholder string
)

var (
	// CreateQueryType create queryType value.
	CreateQueryType queryType = "create"
	// UpdateQueryType update queryType value.
	UpdateQueryType queryType = "update"

	// QuestionMarkQueryPlaceholder question mark queryPlaceholder value.
	QuestionMarkQueryPlaceholder queryPlaceholder = "?"
	// DollarQueryPlaceholder dollar queryPlaceholder value.
	DollarQueryPlaceholder queryPlaceholder = "$"
)

// ErrIsNotStruct if the object given is not a struct.
var ErrIsNotStruct = errors.New("object given is not struct")

func buildQuery(queryType queryType, queryPlaceholder queryPlaceholder, tableName string, columns []string) string {
	var query, parsedColumns, parsedIndexes string

	if queryType == CreateQueryType {
		query = `insert into "` + tableName + `" (%s) values (%s)`
	} else {
		query = `update "` + tableName + `" set (%s) = (%s)`
	}

	for index, column := range columns {
		parsedColumns += `"` + column + `", `
		if queryPlaceholder == DollarQueryPlaceholder {
			parsedIndexes += fmt.Sprintf("$%d, ", (index + 1))
		} else {
			parsedIndexes += "?, "
		}
	}

	if queryType == UpdateQueryType && len(columns) <= 1 {
		query = strings.NewReplacer("(", "", ")", "").Replace(query)
	}
	return fmt.Sprintf(
		query,
		strings.TrimSuffix(parsedColumns, ", "),
		strings.TrimSuffix(parsedIndexes, ", "),
	)
}

// QueryAndValues builds a query (INSERT or UPDATE) by extracting the column
// names of some tag from the struct fields. The column name will be defined
// from the beginning of the tag value until before the first comma found.
//
// The values ​​are taken from the struct fields themselves.
//
// If you do not want any field of the struct to be built together, pass the tag
// value (up to before the first comma) in the "skips" parameter.
//
// The struct passed can be either a real struct or a pointer to a struct.
//
// NOTE: struct fields that are NOT pointers will be automatically discarded
// when building the query.
func QueryAndValues(structure any, queryType queryType, queryPlaceholder queryPlaceholder, tagName, tableName string, skips ...string) (query string, values []any, err error) {
	var columns []string
	valueOfStruct := reflect.ValueOf(structure)
	typeOfStruct := reflect.TypeOf(structure)
	if valueOfStruct.Kind() == reflect.Pointer {
		if valueOfStruct.Elem().Kind() != reflect.Struct {
			return "", nil, ErrIsNotStruct
		}
		valueOfStruct = valueOfStruct.Elem()
		typeOfStruct = typeOfStruct.Elem()

	}
	for index := 0; index < valueOfStruct.NumField(); index++ {
		valueAtualField := valueOfStruct.Field(index)
		typeAtualField := typeOfStruct.Field(index)
		if !valueAtualField.IsValid() || valueAtualField.Kind() != reflect.Pointer || valueAtualField.IsNil() || !typeAtualField.IsExported() {
			continue
		}
		tagValue, hasTag := typeAtualField.Tag.Lookup(tagName)
		parsedTagValue := regexp.MustCompile(`,.*$`).ReplaceAllString(tagValue, "")
		if hasTag && !regexp.MustCompile(parsedTagValue).MatchString(strings.Join(skips, " ")) {
			columns = append(columns, parsedTagValue)
			values = append(values, valueAtualField.Elem().Interface())
		}
	}
	return buildQuery(queryType, queryPlaceholder, tableName, columns), values, nil
}
