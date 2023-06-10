package pgsql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnion1(t *testing.T) {
	sb1 := NewSelectBuilder()
	sb1.Select("id", "name", "created_at")
	sb1.From("demo.user")
	sb1.Where(
		sb1.GreaterThan("id", 1234),
	)

	sb2 := NewSelectBuilder()
	sb2.Select("id", "avatar")
	sb2.From("demo.user_profile")
	sb2.Where(
		sb2.In("status", 1, 2, 5),
	)

	ub := Union(sb1, sb2)
	ub.OrderBy("created_at").Desc()

	result, args := ub.Build()

	assert.Equal(t, "(SELECT id, name, created_at FROM demo.user WHERE id > $1) UNION (SELECT id, avatar FROM demo.user_profile WHERE status IN ($2, $3, $4)) ORDER BY created_at DESC", result)
	assert.Equal(t, []interface{}{1234, 1, 2, 5}, args)
}
