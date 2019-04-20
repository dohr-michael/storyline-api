package arango

import (
	"fmt"
	"github.com/dohrm/go-rsql"
	"strconv"
	"strings"
)

type rsqlParser struct {
	relations []string
	mappings  map[string]string
	variables map[string]interface{}
}

func (p *rsqlParser) expressionToFilter(prefix string, expression rsql.Expression, path string) string {
	switch t := expression.(type) {
	case rsql.AndExpression:
		return p.listToFilter(prefix, t.Items, "and", path)
	case rsql.OrExpression:
		return p.listToFilter(prefix, t.Items, "or", path)
	case rsql.EqualsComparison:
		return p.comparisonToFilter(prefix, t.Identifier, t.Val, "==", path)
	case rsql.NotEqualsComparison:
		return p.comparisonToFilter(prefix, t.Identifier, t.Val, "!=", path)
	case rsql.LikeComparison:
		return p.comparisonToFilter(prefix, t.Identifier, t.Val, "=~", path)
	case rsql.NotLikeComparison:
		return p.comparisonToFilter(prefix, t.Identifier, t.Val, "!~", path)
	case rsql.GreaterThanComparison:
		return p.comparisonToFilter(prefix, t.Identifier, t.Val, ">", path)
	case rsql.GreaterThanOrEqualsComparison:
		return p.comparisonToFilter(prefix, t.Identifier, t.Val, ">=", path)
	case rsql.LessThanComparison:
		return p.comparisonToFilter(prefix, t.Identifier, t.Val, "<", path)
	case rsql.LessThanOrEqualsComparison:
		return p.comparisonToFilter(prefix, t.Identifier, t.Val, "<=", path)
	case rsql.InComparison:
		return p.comparisonToFilter(prefix, t.Identifier, t.Val, "IN", path)
	case rsql.NotInComparison:
		return p.comparisonToFilter(prefix, t.Identifier, t.Val, "NOT IN", path)
	default:
		return ""
	}
}

func (p *rsqlParser) listToFilter(prefix string, expressions []rsql.Expression, separator string, path string) string {
	var list []string
	for idx, expression := range expressions {
		c := p.expressionToFilter(prefix, expression, path+"_"+strconv.FormatInt(int64(idx), 10))
		list = append(list, c)
	}
	return "(" + strings.Join(list, " "+separator+" ") + ")"
}

func (p *rsqlParser) toValue(v rsql.Value) interface{} {
	switch t := v.(type) {
	case rsql.ListValue:
		res := make([]interface{}, 0)
		for _, v := range t.Value {
			res = append(res, p.toValue(v))
		}
		return res
	case rsql.StringValue:
		return t.Value
	case rsql.IntegerValue:
		return t.Value
	case rsql.DoubleValue:
		return t.Value
	case rsql.BooleanValue:
		return t.Value
	case rsql.DateTimeValue:
		return t.Value
	case rsql.DateValue:
		return t.Value
	default:
		return ""
	}
}

func (p *rsqlParser) comparisonToFilter(prefix string, field rsql.Identifier, value rsql.Value, comparator string, path string) string {
	p.variables[path] = p.toValue(value)
	inRelation := false
	fieldName := field.Val
	for _, r := range p.relations {
		if strings.HasPrefix(fieldName, r) {
			inRelation = true
			break
		}
	}
	if p.mappings[fieldName] != "" {
		fieldName = p.mappings[fieldName]
	}
	if inRelation {
		return fmt.Sprintf("%s %s @%s", fieldName, comparator, path)
	}
	return fmt.Sprintf("%s.%s %s @%s", prefix, fieldName, comparator, path)
}

func RsqlToFilter(
	prefix string,
	query rsql.Expression,
	relations []string,
	mappings map[string]string,
) (string, map[string]interface{}) {
	parser := &rsqlParser{
		relations: relations,
		mappings:  mappings,
		variables: make(map[string]interface{}),
	}
	if query == nil {
		return "", parser.variables
	}

	return parser.expressionToFilter(prefix, query, "rsql"), parser.variables
}
