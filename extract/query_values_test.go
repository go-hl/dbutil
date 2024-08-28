package extract

import (
	"fmt"
	"slices"
	"testing"
)

type test struct {
	Column1 *any `query:"column_1,param"`
	Column2 *any `query:"column_2"`
	Column3 any  `query:"column_3,param"`
	Column4 *any `query:"column_4,param"`
	Column5 *any `query:"column_5"`
	Column6 any  `query:"column_6"`
	Column7 *any `query:"column_7,param"`
	Column8 *any `tag:"value_8"`
	Column9 *any `query:"column_9"`
}

type testable struct {
	structure        any
	queryType        queryType
	queryPlaceholder queryPlaceholder
	tagName          string
	tableName        string
	skips            []string

	config       string
	queryExpect  string
	valuesExpect []any
}

var (
	value1 any = "value1"
	value2 any = "value2"
	value4 any = "value4"
	value5 any = "value5"
	value7 any = "value7"
	value8 any = "value8"

	table1 = test{
		Column1: &value1,
		Column2: &value2,
		Column3: "value3",
		Column4: &value4,
		Column5: &value5,
		Column6: "value6",
		Column7: &value7,
		Column8: &value8,
	}
	table2 = test{Column1: &value1}
	table3 test
)

var tests = []testable{
	{table1, CreateQueryType, QuestionMarkQueryPlaceholder, "query", "table", nil, "V1CQZ", `insert into "table" ("column_1", "column_2", "column_4", "column_5", "column_7") values (?, ?, ?, ?, ?)`, []any{"value1", "value2", "value4", "value5", "value7"}},
	{&table1, CreateQueryType, QuestionMarkQueryPlaceholder, "query", "table", nil, "P1CQZ", `insert into "table" ("column_1", "column_2", "column_4", "column_5", "column_7") values (?, ?, ?, ?, ?)`, []any{"value1", "value2", "value4", "value5", "value7"}},
	{table1, CreateQueryType, QuestionMarkQueryPlaceholder, "query", "table", []string{"column_4", "column_7"}, "V1CQS", `insert into "table" ("column_1", "column_2", "column_5") values (?, ?, ?)`, []any{"value1", "value2", "value5"}},
	{table1, CreateQueryType, DollarQueryPlaceholder, "query", "table", nil, "V1CDZ", `insert into "table" ("column_1", "column_2", "column_4", "column_5", "column_7") values ($1, $2, $3, $4, $5)`, []any{"value1", "value2", "value4", "value5", "value7"}},
	{&table1, CreateQueryType, DollarQueryPlaceholder, "query", "table", nil, "P1CDZ", `insert into "table" ("column_1", "column_2", "column_4", "column_5", "column_7") values ($1, $2, $3, $4, $5)`, []any{"value1", "value2", "value4", "value5", "value7"}},
	{table1, CreateQueryType, DollarQueryPlaceholder, "query", "table", []string{"column_4", "column_7"}, "V1CDS", `insert into "table" ("column_1", "column_2", "column_5") values ($1, $2, $3)`, []any{"value1", "value2", "value5"}},
	{table1, UpdateQueryType, QuestionMarkQueryPlaceholder, "query", "table", nil, "V1UQZ", `update "table" set ("column_1", "column_2", "column_4", "column_5", "column_7") = (?, ?, ?, ?, ?)`, []any{"value1", "value2", "value4", "value5", "value7"}},
	{&table1, UpdateQueryType, QuestionMarkQueryPlaceholder, "query", "table", nil, "P1UQZ", `update "table" set ("column_1", "column_2", "column_4", "column_5", "column_7") = (?, ?, ?, ?, ?)`, []any{"value1", "value2", "value4", "value5", "value7"}},
	{table1, UpdateQueryType, QuestionMarkQueryPlaceholder, "query", "table", []string{"column_4", "column_7"}, "V1UQS", `update "table" set ("column_1", "column_2", "column_5") = (?, ?, ?)`, []any{"value1", "value2", "value5"}},
	{table1, UpdateQueryType, DollarQueryPlaceholder, "query", "table", nil, "V1UDZ", `update "table" set ("column_1", "column_2", "column_4", "column_5", "column_7") = ($1, $2, $3, $4, $5)`, []any{"value1", "value2", "value4", "value5", "value7"}},
	{&table1, UpdateQueryType, DollarQueryPlaceholder, "query", "table", nil, "P1UDZ", `update "table" set ("column_1", "column_2", "column_4", "column_5", "column_7") = ($1, $2, $3, $4, $5)`, []any{"value1", "value2", "value4", "value5", "value7"}},
	{table1, UpdateQueryType, DollarQueryPlaceholder, "query", "table", []string{"column_4", "column_7"}, "V1UDS", `update "table" set ("column_1", "column_2", "column_5") = ($1, $2, $3)`, []any{"value1", "value2", "value5"}},
	{table2, CreateQueryType, QuestionMarkQueryPlaceholder, "query", "table", nil, "V2CQZ", `insert into "table" ("column_1") values (?)`, []any{value1}},
	{table2, UpdateQueryType, QuestionMarkQueryPlaceholder, "query", "table", nil, "V2UQZ", `update "table" set "column_1" = ?`, []any{value1}},
	{table3, CreateQueryType, QuestionMarkQueryPlaceholder, "query", "table", nil, "V3CQZ", `insert into "table" () values ()`, []any{}},
	{table3, UpdateQueryType, QuestionMarkQueryPlaceholder, "query", "table", nil, "V3UQZ", `update "table" set  = `, []any{}},
}

func TestQueryAndValues(t *testing.T) {
	for _, test := range tests {
		println("running for:", test.config)
		queryResult, valuesResult, errResult := QueryAndValues(test.structure, test.queryType, test.queryPlaceholder, test.tagName, test.tableName, test.skips...)
		if queryResult != test.queryExpect {
			t.Errorf("\n- CONFIG: %s\n- QUERY result: %s\n- QUERY expect: %s\n", test.config, queryResult, test.queryExpect)
		}
		if !slices.Equal(valuesResult, test.valuesExpect) {
			t.Errorf("\n- CONFIG: %s\n- VALUES result: %#v\n- VALUES expect: %#v\n", test.config, valuesResult, test.valuesExpect)
		}
		if errResult != nil {
			t.Errorf("\n- CONFIG: %s\n- ERROR: %v\n", test.config, errResult)
		}
		fmt.Printf("\t- # query: %s\n\t- $ values: %#v\n\t- ! error: %v\n", queryResult, valuesResult, errResult)
	}
}
