package queryme

import (
	"fmt"
	"testing"
	"github.com/bmizerany/assert"
)

func TestToSql(t *testing.T) {
	p := And{
		Not{Or{
			Lt{"foo", "bob"},
			Eq{"bar", []Value{true}}}},
		Fts{"foo", "go library"}}

	sql, values := ToSql(p, []Field{"foo", "bar"})
	assert.Equal(t, "((NOT (`foo`<? OR `bar`=?)) AND MATCH (`foo`) AGAINST (?))", sql)
	assert.Equal(t, []interface{}{"bob", true, "go library"}, values)

	assert.Panic(t, fmt.Errorf("unauthorized field accessed: %q", "bar"), func() {
		ToSql(p, []Field{"foo"})
	})
}
