package data

import (
	"github.com/dohr-michael/go-libs/filters"
	"github.com/dohrm/go-rsql"
)

func EqStr(field string, value string) filters.Filter {
	return &rsql.EqualsComparison{
		Comparison: rsql.Comparison{
			Identifier: rsql.Identifier{Val: field},
			Val:        rsql.StringValue{Value: value},
		},
	}
}
