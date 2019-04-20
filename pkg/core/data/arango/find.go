package arango

import (
	"bytes"
	"context"
	"github.com/dohr-michael/go-libs/errors"
	"github.com/dohr-michael/go-libs/filters"
	"github.com/dohr-michael/go-libs/storage"
	"github.com/dohr-michael/storyline-api/pkg/core/data"
	"github.com/dohrm/go-rsql"
	"html/template"
	"log"
)

type findParameters struct {
	collection       string
	expression       rsql.Expression
	limit            int64
	offset           int64
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
	RETURN MERGE(c, { {{ range $k, $v := .relationRequests }}{{if $k}}, {{end}}{{$v.fieldName}}: {{$v.fieldName}}{{end}} })
)
RETURN {total: total, items: items}
`

func Find(ctx context.Context, result interface{}, params *findParameters) (*storage.Paged, error) {
	res := &storage.Paged{
		Items: result,
		Query: &filters.Query{
			Filter: params.expression,
			Pager: filters.Pager{
				Limit:  params.limit,
				Offset: params.offset,
			},
		},
	}

	filterStr, args := RsqlToFilter("c",
		params.expression,
		params.relationKeys,
		params.mappings,
	)

	buffer := bytes.Buffer{}
	t, err := template.New("findManyQuery").Parse(findTemplate)
	if err != nil {
		return nil, err
	}
	err = t.Execute(&buffer, map[string]interface{}{
		"collection":       params.collection,
		"offset":           params.offset,
		"limit":            params.limit,
		"relationRequests": params.relationRequests,
		"filter":           filterStr,
	})
	if err != nil {
		return nil, err
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
		return nil, err
	} else if len(items) == 0 {
		return nil, errors.NotFoundError
	}

	doc := items[0]
	s, ok := doc["items"].([]interface{})
	if !ok {
		return nil, errors.NotFoundError
	}
	err = data.Decode(&s, &result)

	if err != nil {
		return nil, err
	}

	res.Total = int64(doc["total"].(float64))

	return res, nil
}
