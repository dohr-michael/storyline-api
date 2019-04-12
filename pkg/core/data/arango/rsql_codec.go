package arango

import (
	"fmt"
	"github.com/dohrm/go-rsql"
	"strconv"
	"strings"
)

func expressionToFilter(prefix string, expression rsql.Expression, path string, params map[string]interface{}) string {
	switch t := expression.(type) {
	case rsql.AndExpression:
		return listToFilter(prefix, t.Items, "and", path, params)
	case rsql.OrExpression:
		return listToFilter(prefix, t.Items, "or", path, params)
	case rsql.EqualsComparison:
		return comparisonToFilter(prefix, t.Identifier, t.Val, "==", path, params)
	case rsql.NotEqualsComparison:
		return comparisonToFilter(prefix, t.Identifier, t.Val, "!=", path, params)
	case rsql.LikeComparison:
		return comparisonToFilter(prefix, t.Identifier, t.Val, "=~", path, params)
	case rsql.NotLikeComparison:
		return comparisonToFilter(prefix, t.Identifier, t.Val, "!~", path, params)
	case rsql.GreaterThanComparison:
		return comparisonToFilter(prefix, t.Identifier, t.Val, ">", path, params)
	case rsql.GreaterThanOrEqualsComparison:
		return comparisonToFilter(prefix, t.Identifier, t.Val, ">=", path, params)
	case rsql.LessThanComparison:
		return comparisonToFilter(prefix, t.Identifier, t.Val, "<", path, params)
	case rsql.LessThanOrEqualsComparison:
		return comparisonToFilter(prefix, t.Identifier, t.Val, "<=", path, params)
	case rsql.InComparison:
		return comparisonToFilter(prefix, t.Identifier, t.Val, "IN", path, params)
	case rsql.NotInComparison:
		return comparisonToFilter(prefix, t.Identifier, t.Val, "NOT IN", path, params)
	default:
		return ""
	}
}

func listToFilter(prefix string, expressions []rsql.Expression, separator string, path string, params map[string]interface{}) string {
	var list []string
	for idx, expression := range expressions {
		c := expressionToFilter(prefix, expression, path+"_"+strconv.FormatInt(int64(idx), 10), params)
		list = append(list, c)
	}
	return "(" + strings.Join(list, " "+separator+" ") + ")"
}

func toValue(v rsql.Value) interface{} {
	switch t := v.(type) {
	case rsql.ListValue:
		res := make([]interface{}, 0)
		for _, v := range t.Value {
			res = append(res, toValue(v))
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

func comparisonToFilter(prefix string, field rsql.Identifier, value rsql.Value, comparator string, path string, params map[string]interface{}) string {
	params[path] = toValue(value)
	return fmt.Sprintf("%s.%s %s @%s", prefix, field.Val, comparator, path)
}

func RsqlToFilter(prefix string, query rsql.Expression) (string, map[string]interface{}) {
	res := map[string]interface{}{}
	if query == nil {
		return "", res
	}
	return expressionToFilter(prefix, query, "rsql", res), res
}
