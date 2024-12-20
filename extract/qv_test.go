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
	returning        string
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
	{table1, CreateType, QuestionMarkPlaceholder, "query", "table", "", nil, "V1CQZN", `insert into "table" ("column_1", "column_2", "column_4", "column_5", "column_7") values (?, ?, ?, ?, ?)`, []any{"value1", "value2", "value4", "value5", "value7"}},
	{&table1, CreateType, QuestionMarkPlaceholder, "query", "table", "", nil, "P1CQZN", `insert into "table" ("column_1", "column_2", "column_4", "column_5", "column_7") values (?, ?, ?, ?, ?)`, []any{"value1", "value2", "value4", "value5", "value7"}},
	{table1, CreateType, QuestionMarkPlaceholder, "query", "table", "", []string{"column_4", "column_7"}, "V1CQSN", `insert into "table" ("column_1", "column_2", "column_5") values (?, ?, ?)`, []any{"value1", "value2", "value5"}},
	{table1, CreateType, DollarPlaceholder, "query", "table", "", nil, "V1CDZN", `insert into "table" ("column_1", "column_2", "column_4", "column_5", "column_7") values ($1, $2, $3, $4, $5)`, []any{"value1", "value2", "value4", "value5", "value7"}},
	{&table1, CreateType, DollarPlaceholder, "query", "table", "", nil, "P1CDZN", `insert into "table" ("column_1", "column_2", "column_4", "column_5", "column_7") values ($1, $2, $3, $4, $5)`, []any{"value1", "value2", "value4", "value5", "value7"}},
	{table1, CreateType, DollarPlaceholder, "query", "table", "", []string{"column_4", "column_7"}, "V1CDSN", `insert into "table" ("column_1", "column_2", "column_5") values ($1, $2, $3)`, []any{"value1", "value2", "value5"}},
	{table1, UpdateType, QuestionMarkPlaceholder, "query", "table", "", nil, "V1UQZN", `update "table" set ("column_1", "column_2", "column_4", "column_5", "column_7") = (?, ?, ?, ?, ?)`, []any{"value1", "value2", "value4", "value5", "value7"}},
	{&table1, UpdateType, QuestionMarkPlaceholder, "query", "table", "", nil, "P1UQZN", `update "table" set ("column_1", "column_2", "column_4", "column_5", "column_7") = (?, ?, ?, ?, ?)`, []any{"value1", "value2", "value4", "value5", "value7"}},
	{table1, UpdateType, QuestionMarkPlaceholder, "query", "table", "", []string{"column_4", "column_7"}, "V1UQSN", `update "table" set ("column_1", "column_2", "column_5") = (?, ?, ?)`, []any{"value1", "value2", "value5"}},
	{table1, UpdateType, DollarPlaceholder, "query", "table", "", nil, "V1UDZN", `update "table" set ("column_1", "column_2", "column_4", "column_5", "column_7") = ($1, $2, $3, $4, $5)`, []any{"value1", "value2", "value4", "value5", "value7"}},
	{&table1, UpdateType, DollarPlaceholder, "query", "table", "", nil, "P1UDZN", `update "table" set ("column_1", "column_2", "column_4", "column_5", "column_7") = ($1, $2, $3, $4, $5)`, []any{"value1", "value2", "value4", "value5", "value7"}},
	{table1, UpdateType, DollarPlaceholder, "query", "table", "", []string{"column_4", "column_7"}, "V1UDSN", `update "table" set ("column_1", "column_2", "column_5") = ($1, $2, $3)`, []any{"value1", "value2", "value5"}},
	{table2, CreateType, QuestionMarkPlaceholder, "query", "table", "", nil, "V2CQZN", `insert into "table" ("column_1") values (?)`, []any{value1}},
	{table2, UpdateType, QuestionMarkPlaceholder, "query", "table", "", nil, "V2UQZN", `update "table" set "column_1" = ?`, []any{value1}},
	{table3, CreateType, QuestionMarkPlaceholder, "query", "table", "", nil, "V3CQZN", `insert into "table" () values ()`, []any{}},
	{table3, UpdateType, QuestionMarkPlaceholder, "query", "table", "", nil, "V3UQZN", `update "table" set  = `, []any{}},
	{table1, CreateType, DollarPlaceholder, "query", "table", "column_x", []string{"column_4", "column_7"}, "V1CQSR", `insert into "table" ("column_1", "column_2", "column_5") values ($1, $2, $3) returning "column_x"`, []any{"value1", "value2", "value5"}},
	{table2, CreateType, DollarPlaceholder, "query", "table", "column_x", nil, "V2CQZR", `insert into "table" ("column_1") values ($1) returning "column_x"`, []any{value1}},
	{table3, CreateType, DollarPlaceholder, "query", "table", "column_x", nil, "V3CQZR", `insert into "table" () values () returning "column_x"`, []any{}},
}

func TestQueryAndValues(t *testing.T) {
	for _, test := range tests {
		println("running for:", test.config)
		queryResult, valuesResult, errResult := QueryAndValues(test.structure, test.queryType, test.queryPlaceholder, test.tagName, test.tableName, test.returning, test.skips...)
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
