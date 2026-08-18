package sql_parser

import (
	"fmt"
	"strings"

	"github.com/openshift-online/ocm-common/pkg/utils/parser/state_machine"
	"github.com/openshift-online/ocm-common/pkg/utils/parser/string_parser"
)

const defaultMaximumComplexity = 10

// SQLParser - This object is to be used to parse and validate WHERE clauses (only portion after the `WHERE` is supported)
type SQLParser interface {
	// Parse - parses the received SQL string and returns the parsed values or an error
	// Returns:
	// - string: The parsed SQL replacing all the values with '?' placeholders
	// - interface{}: All the values to pass to the database (to replace the '?' placeholders)
	// - error: non nil in case of any error
	Parse(sql string) (string, interface{}, error)
}

type sqlParser struct {
	// configuration
	maximumComplexity int
	parser            *string_parser.StringParser

	// current parsing state
	// counts the number of joins
	complexity int
	// counts the number of braces to be closed
	openBraces   int
	validColumns []string
	columnPrefix string

	// current parsing result
	resultQry    string
	resultValues []interface{}
}

var _ SQLParser = &sqlParser{}

func (p *sqlParser) Parse(sql string) (string, interface{}, error) {
	p.reset()

	if err := p.parser.Parse(sql); err != nil {
		return "", nil, err
	}

	if p.openBraces > 0 {
		return "", nil, fmt.Errorf("EOF while searching for closing brace ')'")
	}

	p.resultQry = strings.Trim(p.resultQry, " ")
	return p.resultQry, p.resultValues, nil
}

func (p *sqlParser) reset() {
	p.complexity = 0
	p.openBraces = 0
	p.resultQry = ""
	p.resultValues = nil
}

func (p *sqlParser) transitionInterceptor(_, to *state_machine.State[string, string], tokenValue string) error {
	countOpenBraces := func(tok string) error {
		switch tok {
		case "(":
			p.openBraces++
		case ")":
			p.openBraces--
		}
		if p.openBraces < 0 {
			return fmt.Errorf("unexpected ')'")
		}
		return nil
	}

	tokenFamily := to.Data() // The grammar configures the custom state data as the token family
	switch tokenFamily {
	case braceTokenFamily:
		if err := countOpenBraces(tokenValue); err != nil {
			return err
		}
		p.resultQry += tokenValue
		return nil
	case valueTokenFamily:
		p.resultQry += " ?"
		p.resultValues = append(p.resultValues, tokenValue)
		return nil
	case quotedValueTokenFamily:
		p.resultQry += " ?"
		p.resultValues = append(p.resultValues, unquoteValue(tokenValue))
		return nil
	case logicalOpTokenFamily:
		p.complexity++
		if p.complexity > p.maximumComplexity {
			return fmt.Errorf("maximum number of permitted joins (%d) exceeded", p.maximumComplexity)
		}
		p.resultQry += " " + tokenValue + " "
		return nil
	case columnTokenFamily:
		// we want column names to be lowercase
		columnName := strings.ToLower(tokenValue)
		if len(p.validColumns) > 0 && !contains(p.validColumns, columnName) {
			return fmt.Errorf("invalid column name: '%s', valid values are: %v", tokenValue, p.validColumns)
		}
		if p.columnPrefix != "" && !strings.HasPrefix(columnName, p.columnPrefix+".") {
			columnName = p.columnPrefix + "." + columnName
		}
		p.resultQry += columnName
		return nil
	case jsonbFamily:
		// The JSONB operators (`->`, `->>`, `@>`) are fixed grammar literals and are
		// safe to emit verbatim. The JSONB path keys (jsonbField / jsonbFieldToStringify)
		// are user supplied and MUST be parameterised, otherwise a quoted key can break
		// out of the string literal and inject arbitrary SQL.
		if to.Name() == jsonbField || to.Name() == jsonbFieldToStringify {
			p.resultQry += " ?"
			p.resultValues = append(p.resultValues, unquoteValue(tokenValue))
		} else {
			p.resultQry += " " + tokenValue
		}
		return nil
	default:
		p.resultQry += " " + tokenValue
		return nil
	}
}

func unquoteValue(tokenValue string) string {
	// unescape
	tmp := strings.ReplaceAll(tokenValue, `\'`, "'")
	// remove quotes:
	if len(tmp) > 1 {
		tmp = string([]rune(tmp)[1 : len(tmp)-1])
	}
	return tmp
}

func contains(ary []string, value string) bool {
	for _, v := range ary {
		if v == value {
			return true
		}
	}
	return false
}
