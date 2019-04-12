package arango

import (
	"encoding/json"
	"time"
)

type Identifier struct {
	Key    string `json:"_key" mapstructure:"_key"`
	Rev    string `json:"_rev" mapstructure:"_rev"`
	NodeId string `json:"_id" mapstructure:"_id"`
}

type Edge struct {
	From string `json:"_from" mapstructure:"_from"`
	To   string `json:"_to" mapstructure:"_to"`
}

type Time time.Time

func (t Time) MarshalJSON() ([]byte, error) {
	f := time.Now().UnixNano() / int64(time.Millisecond)
	return json.Marshal(f)
}
