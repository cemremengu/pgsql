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

	sb2 := newSelectBuilder()
	sb2.Select("id", "avatar")
	sb2.From("demo.user_profile")
	sb2.Where(
		sb2.In("status", 1, 2, 5),
	)

	ub := Union(sb1, sb2)
	ub.OrderBy("created_at").Desc()

	result, args := ub.Build()

	assert.Equal(t, "UPDATE demo.user SET visited = visited + 1 WHERE id = 1234", result)
	assert.Empty(t, args)
}
