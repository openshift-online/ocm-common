package sql_parser

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SQLParser", func() {
	type testData struct {
		qry        string
		outQry     string
		outValues  []interface{}
		wantErr    bool
		errMessage string
	}

	parserTest := func(tt testData, parser SQLParser) {
		qry, values, err := parser.Parse(tt.qry)

		if !tt.wantErr {
			Expect(err).ToNot(HaveOccurred(), "QueryParser() error = %v, wantErr = %v", err, tt.wantErr)
		}

		Expect(err != nil).To(Equal(tt.wantErr))

		if err == nil && tt.outQry != "" {
			if tt.outQry != "" {
				Expect(qry).To(Equal(tt.outQry))
			}
			if tt.outValues != nil {
				Expect(values).To(Equal(tt.outValues))
			}
		}

		if err != nil && tt.wantErr && tt.errMessage != "" {
			Expect(err.Error()).To(Equal(tt.errMessage))
		}
	}

	DescribeTable("Basic Parsing", parserTest,
		Entry("Just `=` sign", testData{
			qry:        "=",
			wantErr:    true,
			errMessage: "[1] error parsing the filter: unexpected token `=`",
		}, NewSQLParser()),
		Entry("Incomplete query", testData{
			qry:        "name=",
			wantErr:    true,
			errMessage: "EOF encountered while parsing string",
		}, NewSQLParser()),
		Entry("Incomplete join", testData{
			qry:        "name='test' and ",
			wantErr:    true,
			errMessage: "EOF encountered while parsing string",
		}, NewSQLParser()),
		Entry("Escaped quote", testData{
			qry:       `name='test\'123'`,
			outQry:    "name = ?",
			outValues: []interface{}{"test'123"},
			wantErr:   false,
		}, NewSQLParser()),
		Entry("Wrong unescaped quote", testData{
			qry:        `name='test'123'`,
			wantErr:    true,
			errMessage: "[12] error parsing the filter: unexpected token `123`",
		}, NewSQLParser()),
		Entry("Quoted parenthesis", testData{
			qry:       `name='test(123)'`,
			wantErr:   false,
			outQry:    "name = ?",
			outValues: []interface{}{"test(123)"},
		}, NewSQLParser()),
		Entry("Quoted special characters", testData{
			qry:       `name='@,\\'""(){}/'`,
			wantErr:   false,
			outQry:    "name = ?",
			outValues: []interface{}{`@,\'""(){}/`},
		}, NewSQLParser()),
	)

	DescribeTable("IN Keyword Parsing", parserTest,
		Entry("IN keyword", testData{
			qry:       "name IN ('value1', 'value2')",
			outQry:    "name IN( ? , ?)",
			outValues: []interface{}{"value1", "value2"},
			wantErr:   false,
		}, NewSQLParser()),
		Entry("IN with single value", testData{
			qry:       "name IN ('value1')",
			outQry:    "name IN( ?)",
			outValues: []interface{}{"value1"},
			wantErr:   false,
		}, NewSQLParser()),
		Entry("IN with no values", testData{
			qry:        "name IN ()",
			outQry:     "",
			outValues:  nil,
			wantErr:    true,
			errMessage: "[10] error parsing the filter: unexpected token `)`",
		}, NewSQLParser()),
		Entry("invalid IN (ends with comma)", testData{
			qry:        "name IN ('value1',)",
			outQry:     "",
			outValues:  nil,
			wantErr:    true,
			errMessage: "[19] error parsing the filter: unexpected token `)`",
		}, NewSQLParser()),
		Entry("invalid IN (no closed brace)", testData{
			qry:        "name IN ('value1'",
			outQry:     "",
			outValues:  nil,
			wantErr:    true,
			errMessage: "EOF encountered while parsing string",
		}, NewSQLParser()),
		Entry("IN in complex query", testData{
			qry:       "((cloud_provider = Value and name = value1) and (owner <> value2 or region=b ) or owner in ('owner1', 'owner2', 'owner3')) or owner=c or name=e and region LIKE '%test%' and instance_type=standard",
			outQry:    "((cloud_provider = ? and name = ?) and (owner <> ? or region = ?) or owner in( ? , ? , ?)) or owner = ? or name = ? and region LIKE ? and instance_type = ?",
			outValues: []interface{}{"Value", "value1", "value2", "b", "owner1", "owner2", "owner3", "c", "e", "%test%", "standard"},
			wantErr:   false,
		}, NewSQLParser()),
		Entry("IN with non quoted and quoted values", testData{
			qry:       "owner in (owner1, 'owner2', owner3)",
			outQry:    "owner in( ? , ? , ?)",
			outValues: []interface{}{"owner1", "owner2", "owner3"},
			wantErr:   false,
		}, NewSQLParser()),
		Entry("IN with quoted value containing a comma", testData{
			qry:       "owner in (owner1, 'owner2,', owner3)",
			outQry:    "owner in( ? , ? , ?)",
			outValues: []interface{}{"owner1", "owner2,", "owner3"},
			wantErr:   false,
		}, NewSQLParser()),
		Entry("negated IN in complex query", testData{
			qry:       "((cloud_provider = Value and name = value1) and (owner <> value2 or region=b ) or owner not in ('owner1', 'owner2', 'owner3')) or owner=c or name=e and region LIKE '%test%'",
			outQry:    "((cloud_provider = ? and name = ?) and (owner <> ? or region = ?) or owner not  in( ? , ? , ?)) or owner = ? or name = ? and region LIKE ?",
			outValues: []interface{}{"Value", "value1", "value2", "b", "owner1", "owner2", "owner3", "c", "e", "%test%"},
			wantErr:   false,
		}, NewSQLParser()),
	)

	DescribeTable("BRACES validation", parserTest,
		Entry("Complex query with braces", testData{
			qry:       "((cloud_provider = Value and name = value1) and (owner <> value2 or region=b ) ) or owner=c or name=e and region LIKE '%test%'",
			outQry:    "((cloud_provider = ? and name = ?) and (owner <> ? or region = ?)) or owner = ? or name = ? and region LIKE ?",
			outValues: []interface{}{"Value", "value1", "value2", "b", "c", "e", "%test%"},
			wantErr:   false,
		}, NewSQLParser()),
		Entry("Complex query with braces and quoted values with escaped quote", testData{
			qry:       `((cloud_provider = 'Value' and name = 'val\'ue1') and (owner = value2 or region='b' ) ) or owner=c or name=e and region LIKE '%test%'`,
			outQry:    "((cloud_provider = ? and name = ?) and (owner = ? or region = ?)) or owner = ? or name = ? and region LIKE ?",
			outValues: []interface{}{"Value", "val'ue1", "value2", "b", "c", "e", "%test%"},
			wantErr:   false,
		}, NewSQLParser()),
		Entry("Complex query with braces and quoted values with spaces", testData{
			qry:       `((cloud_provider = 'Value' and name = 'val ue1') and (owner = ' value2  ' or region='b' ) ) or owner=c or name=e and region LIKE '%test%'`,
			outQry:    "((cloud_provider = ? and name = ?) and (owner = ? or region = ?)) or owner = ? or name = ? and region LIKE ?",
			outValues: []interface{}{"Value", "val ue1", " value2  ", "b", "c", "e", "%test%"},
			wantErr:   false,
		}, NewSQLParser()),
		Entry("Complex query with braces and empty quoted values", testData{
			qry:       `((cloud_provider = 'Value' and name = '') and (owner = ' value2  ' or region='' ) ) or owner=c or name=e and region LIKE '%test%'`,
			outQry:    "((cloud_provider = ? and name = ?) and (owner = ? or region = ?)) or owner = ? or name = ? and region LIKE ?",
			outValues: []interface{}{"Value", "", " value2  ", "", "c", "e", "%test%"},
			wantErr:   false,
		}, NewSQLParser()),
	)

	DescribeTable("ILIKE Keyword Parsing", parserTest,
		Entry("Complex query with braces", testData{
			qry:       "((cloud_provider = Value and name = value1) and (owner <> value2 or region=b ) ) or owner=c or name=e and region LIKE '%test%'",
			outQry:    "((cloud_provider = ? and name = ?) and (owner <> ? or region = ?)) or owner = ? or name = ? and region LIKE ?",
			outValues: []interface{}{"Value", "value1", "value2", "b", "c", "e", "%test%"},
			wantErr:   false,
		}, NewSQLParser()),
	)

	DescribeTable("JSONB Query Parsing", parserTest,
		Entry("JSONB query", testData{
			qry:       `manifest->'data'->'manifest'->'metadata'->'labels'->>'foo' = 'bar'`,
			outQry:    "manifest -> ? -> ? -> ? -> ? ->> ? = ?",
			outValues: []interface{}{"data", "manifest", "metadata", "labels", "foo", "bar"},
			wantErr:   false,
		}, NewSQLParser()),
		Entry("Invalid JSONB query", testData{
			qry:        `manifest->'data'->'manifest'->'metadata'->'labels'->'foo' = 'bar'`,
			outQry:     "manifest -> 'data' -> 'manifest' -> 'metadata' -> 'labels' ->> 'foo' = ?",
			outValues:  nil,
			wantErr:    true,
			errMessage: "[59] error parsing the filter: unexpected token `=`",
		}, NewSQLParser()),
		Entry("Complex JSONB query", testData{
			qry: `manifest->'data'->'manifest'->'metadata'->'labels'->>'foo' = 'bar' and ` +
				`( manifest->'data'->'manifest' ->> 'foo' in ('value1', 'value2') or ` +
				`manifest->'data'->'manifest'->>'labels' <> 'foo1')`,
			outQry: "manifest -> ? -> ? -> ? -> ? ->> ? = ? and " +
				"(manifest -> ? -> ? ->> ? in( ? , ?) or " +
				"manifest -> ? -> ? ->> ? <> ?)",
			outValues: []interface{}{
				"data", "manifest", "metadata", "labels", "foo", "bar",
				"data", "manifest", "foo", "value1", "value2",
				"data", "manifest", "labels", "foo1"},
			wantErr: false,
		}, NewSQLParser()),
		Entry("JSONB Query @>", testData{
			qry:       `resources.payload -> 'data' -> 'manifests' @> '[{"metadata":{"labels":{"foo":"bar"}}}]'`,
			outQry:    "resources.payload -> ? -> ? @> ?",
			outValues: []interface{}{"data", "manifests", `[{"metadata":{"labels":{"foo":"bar"}}}]`},
			wantErr:   false,
		}, NewSQLParser()),
		Entry("Mixed JSONB Query", testData{
			qry: `manifest->'data'->'manifest'->'metadata'->'labels'->>'foo' = 'bar' and ` +
				`( manifest->'data'->'manifest' ->> 'foo' in ('value1', 'value2') or ` +
				`manifest->'data'->'manifest'->>'labels' <> 'foo1')` +
				` AND resources.payload -> 'data' -> 'manifests' @> '[{"metadata":{"labels":{"foo":"bar"}}}]' OR ` +
				` my_column in (1, 2, 3) and my_column2 = 'value'`,
			outQry: "manifest -> ? -> ? -> ? -> ? ->> ? = ? " +
				"and (manifest -> ? -> ? ->> ? in( ? , ?) " +
				"or manifest -> ? -> ? ->> ? <> ?) " +
				"AND resources.payload -> ? -> ? @> ? " +
				"OR my_column in( ? , ? , ?) and my_column2 = ?",
			outValues: []interface{}{
				"data", "manifest", "metadata", "labels", "foo", "bar",
				"data", "manifest", "foo", "value1", "value2",
				"data", "manifest", "labels", "foo1",
				"data", "manifests", `[{"metadata":{"labels":{"foo":"bar"}}}]`,
				"1", "2", "3", "value"},
			wantErr: false,
		}, NewSQLParser()),
		Entry("JSONB key with injection payload is parameterised", testData{
			// A malicious JSONB key that, under the previous verbatim concatenation,
			// would break out of the string literal and inject SQL. It must now be
			// emitted as a bind parameter instead of reaching the query text.
			qry:       `manifest->'x\' OR \'1\'=\'1'->>'foo' = 'bar'`,
			outQry:    "manifest -> ? ->> ? = ?",
			outValues: []interface{}{`x' OR '1'='1`, "foo", "bar"},
			wantErr:   false,
		}, NewSQLParser()),
	)

	DescribeTable("MAXIMUM COMPLEXITY", parserTest,
		Entry("Complexity ok", testData{
			qry:       "((cloud_provider = Value and name = value1) and (owner <> value2 or region=b ) ) or owner=c or name=e and region LIKE '%test%'",
			outQry:    "((cloud_provider = ? and name = ?) and (owner <> ? or region = ?)) or owner = ? or name = ? and region LIKE ?",
			outValues: []interface{}{"Value", "value1", "value2", "b", "c", "e", "%test%"},
			wantErr:   false,
		}, NewSQLParser()),
		Entry("MaximumComplexity exceeded", testData{
			qry:        "((cloud_provider = Value and name = value1) and (owner <> value2 or region=b ) ) or owner=c or name=e and region LIKE '%test%'",
			wantErr:    true,
			errMessage: "[82] error parsing the filter: maximum number of permitted joins (3) exceeded",
		}, NewSQLParser(WithMaximumComplexity(3))),
	)

	DescribeTable("ALLOWED COLUMNS", parserTest,
		Entry("Any Column", testData{
			qry:       "((cloud_provider = Value and name = value1) and (owner <> value2 or region=b ) ) or owner=c or name=e and region LIKE '%test%'",
			outQry:    "((cloud_provider = ? and name = ?) and (owner <> ? or region = ?)) or owner = ? or name = ? and region LIKE ?",
			outValues: []interface{}{"Value", "value1", "value2", "b", "c", "e", "%test%"},
			wantErr:   false,
		}, NewSQLParser()),
		Entry("Enlisted columns - ok", testData{
			qry:       "((cloud_provider = Value and name = value1) and (owner <> value2 or region=b ) ) or owner=c or name=e and region LIKE '%test%'",
			outQry:    "((cloud_provider = ? and name = ?) and (owner <> ? or region = ?)) or owner = ? or name = ? and region LIKE ?",
			outValues: []interface{}{"Value", "value1", "value2", "b", "c", "e", "%test%"},
			wantErr:   false,
		}, NewSQLParser(WithValidColumns("cloud_provider", "name", "owner", "region"))),
		Entry("Enlisted columns - fail", testData{
			qry:        "((cloud_provider = Value and name = value1) and (owner <> value2 or region=b ) ) or owner=c or name=e and region LIKE '%test%'",
			outQry:     "((cloud_provider = ? and name = ?) and (owner <> ? or region = ?)) or owner = ? or name = ? and region LIKE ?",
			outValues:  []interface{}{"Value", "value1", "value2", "b", "c", "e", "%test%"},
			wantErr:    true,
			errMessage: "[50] error parsing the filter: invalid column name: 'owner', valid values are: [cloud_provider name region]",
		}, NewSQLParser(WithValidColumns("cloud_provider", "name", "region"))),
	)

	DescribeTable("COLUMN PREFIX", parserTest,
		Entry("Empty prefix", testData{
			qry:       "((cloud_provider = Value and name = value1) and (owner <> value2 or region=b ) ) or owner=c or name=e and region LIKE '%test%'",
			outQry:    "((cloud_provider = ? and name = ?) and (owner <> ? or region = ?)) or owner = ? or name = ? and region LIKE ?",
			outValues: []interface{}{"Value", "value1", "value2", "b", "c", "e", "%test%"},
			wantErr:   false,
		}, NewSQLParser()),
		Entry("All spaces prefix", testData{
			qry:       "((cloud_provider = Value and name = value1) and (owner <> value2 or region=b ) ) or owner=c or name=e and region LIKE '%test%'",
			outQry:    "((cloud_provider = ? and name = ?) and (owner <> ? or region = ?)) or owner = ? or name = ? and region LIKE ?",
			outValues: []interface{}{"Value", "value1", "value2", "b", "c", "e", "%test%"},
			wantErr:   false,
		}, NewSQLParser(WithColumnPrefix("   "))),
		Entry("custom prefix", testData{
			qry:       "((cloud_provider = Value and name = value1) and (owner <> value2 or region=b ) ) or owner=c or name=e and region LIKE '%test%'",
			outQry:    "((main.cloud_provider = ? and main.name = ?) and (main.owner <> ? or main.region = ?)) or main.owner = ? or main.name = ? and main.region LIKE ?",
			outValues: []interface{}{"Value", "value1", "value2", "b", "c", "e", "%test%"},
			wantErr:   false,
		}, NewSQLParser(WithColumnPrefix("main"))),
	)
})
