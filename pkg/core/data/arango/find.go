package arango

import (
	"bytes"
	"context"
	"fmt"
	"github.com/dohr-michael/go-libs/errors"
	"github.com/dohr-michael/go-libs/filters"
	"github.com/dohr-michael/storyline-api/pkg/core/data"
	"github.com/dohrm/go-rsql"
	"html/template"
	"log"
	"strings"
)

type sortDirection string

const Asc = sortDirection("ASC")
const Desc = sortDirection("DESC")

type sorter struct {
	by        string
	direction sortDirection
}

type findParameters struct {
	collection       string
	expression       rsql.Expression
	limit            int64
	offset           int64
	sortBy           []sorter
	relationRequests []map[string]interface{}
	mappings         map[string]string
	relationKeys     []string
}

func NewFind(collection string) *findParameters {
	return &findParameters{
		collection: collection,
		offset:     0,
		limit:      100,
		mappings: map[string]string{
			"_key": "_key",
		},
		sortBy:           []sorter{},
		relationRequests: []map[string]interface{}{},
	}
}

func (p *findParameters) WithPager(offset, limit int64) *findParameters {
	p.offset = offset
	p.limit = limit
	return p
}

func (p *findParameters) WithFieldId(fieldId string) *findParameters {
	delete(p.mappings, "_key")
	p.mappings[fieldId] = "_key"
	return p
}

func (p *findParameters) WithRsql(filter string) (*findParameters, error) {
	exp, err := rsql.Parse(filter)
	if err != nil {
		return p, err
	}
	p.expression = exp
	return p, err
}

func (p *findParameters) WithQuery(query *filters.Query) *findParameters {
	p.limit = query.Pager.Limit
	p.offset = query.Pager.Offset
	p.expression = query.Filter
	return p
}

func (p *findParameters) WithRelationRequests(request ...*relationRequest) *findParameters {
	for _, v := range request {
		p.relationRequests = append(p.relationRequests, map[string]interface{}{
			"collection": v.collection,
			"direction":  string(v.direction),
			"fieldName":  v.fieldName,
			"kind":       v.kind,
		})
		p.relationKeys = append(p.relationKeys, v.fieldName)
		for modelId, dataId := range v.relationFieldMapping {
			p.mappings[v.fieldName+"."+modelId] = v.fieldName + "." + dataId
		}
	}
	return p
}

func (p *findParameters) WithSortByAsc(by string) *findParameters {
	p.sortBy = append(p.sortBy, sorter{by: by, direction: Asc})
	return p
}

func (p *findParameters) WithSortByDesc(by string) *findParameters {
	p.sortBy = append(p.sortBy, sorter{by: by, direction: Desc})
	return p
}

const findTemplate = `
let total = LENGTH(
FOR c IN {{.collection}}
	{{ range $k, $v := .relationRequests }}FOR {{$v.fieldName}}, edge{{$k}} IN 1..1 {{$v.direction}} c {{$v.collection}} FILTER edge{{$k}}.kind == "{{$v.kind}}"
	{{end}}{{if .filter}}FILTER {{.filter}}
	{{end}}RETURN 1
)
let items = (
FOR c IN {{.collection}}
	{{ range $k, $v := .relationRequests }}FOR {{$v.fieldName}}, edge{{$k}} IN 1..1 {{$v.direction}} c {{$v.collection}} FILTER edge{{$k}}.kind == "{{$v.kind}}"
	{{end}}{{if .filter}}FILTER {{.filter}}
	{{end}}LIMIT {{.offset}}, {{.limit}}
	{{ if .sortBy }}SORT {{ .sortBy }}{{end}}
	RETURN MERGE(c, { {{ range $k, $v := .relationRequests }}{{if $k}}, {{end}}{{$v.fieldName}}: {{$v.fieldName}}{{end}} })
)
RETURN {total: total, items: items}
`

func Find(ctx context.Context, result interface{}, params *findParameters) (int64, error) {
	const prefix = "c"

	filterStr, args := RsqlToFilter(prefix,
		params.expression,
		params.relationKeys,
		params.mappings,
	)

	buffer := bytes.Buffer{}
	t, err := template.New("findManyQuery").Parse(findTemplate)
	if err != nil {
		return 0, err
	}
	var sortBy []string
	for _, s := range params.sortBy {
		sortBy = append(sortBy, fmt.Sprintf("%s.%s %s", prefix, s.by, s.direction))
	}

	err = t.Execute(&buffer, map[string]interface{}{
		"collection":       params.collection,
		"offset":           params.offset,
		"limit":            params.limit,
		"relationRequests": params.relationRequests,
		"filter":           filterStr,
		"sortBy":           strings.Join(sortBy, ","),
	})
	if err != nil {
		return 0, err
	}
	for k, v := range args {
		log.Printf("%s = %v", k, v)
	}
	log.Printf(buffer.String())

	items, err := RunQuery(
		ctx,
		NewRunQuery(buffer.String()).WithParams(args),
	)
	if err != nil {
		return 0, err
	} else if len(items) == 0 {
		return 0, errors.NotFoundError
	}

	doc := items[0]
	s, ok := doc["items"].([]interface{})
	if !ok {
		return 0, errors.NotFoundError
	}
	err = data.Decode(&s, &result)

	if err != nil {
		return 0, err
	}

	return int64(doc["total"].(float64)), nil
}
