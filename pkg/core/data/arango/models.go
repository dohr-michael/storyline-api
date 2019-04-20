package arango

type QueryResult map[string]interface{}
type QueryResults []QueryResult

type Identifier struct {
	Key    string `json:"_key" mapstructure:"_key"`
	Rev    string `json:"_rev" mapstructure:"_rev"`
	NodeId string `json:"_id" mapstructure:"_id"`
}

type Edge struct {
	From string `json:"_from" mapstructure:"_from"`
	To   string `json:"_to" mapstructure:"_to"`
}

type relationDirection string

const (
	InDirection  = relationDirection("INBOUND")
	OutDirection = relationDirection("OUTBOUND")
)

type relationRequest struct {
	fieldName            string
	kind                 string
	collection           string
	direction            relationDirection
	relationFieldMapping map[string]string
}

func NewRelationRequest(
	relationCollection string,
	fieldName string,
	direction relationDirection,
	kind string,
) *relationRequest {
	return &relationRequest{
		fieldName:            fieldName,
		direction:            direction,
		kind:                 kind,
		collection:           relationCollection,
		relationFieldMapping: map[string]string{},
	}
}

func (r *relationRequest) WithFieldMapping(modelId string, dataId string) *relationRequest {
	r.relationFieldMapping[modelId] = dataId
	return r
}
