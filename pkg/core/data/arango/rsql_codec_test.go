package arango

import (
	"github.com/dohrm/go-rsql"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestRsqlToFilter(t *testing.T) {
	toTest1, _ := rsql.Parse(`(a==12 or a==13) and (b=="42" or b==false) and d.id=="666"`)
	res, params := RsqlToFilter("t",
		toTest1,
		[]string{"d"},
		map[string]string{"d.id": "d._key"},
	)
	assert.Equal(t, res, "((t.a == @rsql_0_0 or t.a == @rsql_0_1) and (t.b == @rsql_1_0 or t.b == @rsql_1_1) and d._key == @rsql_2)")
	assert.Equal(t, params["rsql_0_0"], int64(12))
	assert.Equal(t, params["rsql_0_1"], int64(13))
	assert.Equal(t, params["rsql_1_0"], "42")
	assert.Equal(t, params["rsql_1_1"], false)
	assert.Equal(t, params["rsql_2"], "666")
}
