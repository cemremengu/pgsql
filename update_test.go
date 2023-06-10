package pgsql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdate1(t *testing.T) {
	result, args := Update("demo.user").
		Set(
			"visited = visited + 1",
		).
		Where(
			"id = 1234",
		).
		Build()

	assert.Equal(t, "UPDATE demo.user SET visited = visited + 1 WHERE id = 1234", result)
	assert.Empty(t, args)
}

func TestUpdate2(t *testing.T) {
	ub := NewUpdateBuilder()
	ub.Update("demo.user")
	ub.Set(
		ub.Assign("type", "sys"),
		ub.Incr("credit"),
		"modified_at = UNIX_TIMESTAMP(NOW())", // It's allowed to write arbitrary SQL.
	)
	ub.Where(
		ub.GreaterThan("id", 1234),
		ub.Like("name", "%Du"),
		ub.Or(
			ub.IsNull("id_card"),
			ub.In("status", 1, 2, 5),
		),
		"modified_at > created_at + "+ub.Var(86400), // It's allowed to write arbitrary SQL.
	)
	ub.OrderBy("id").Asc()

	result, args := ub.Build()

	assert.Equal(t, "UPDATE demo.user SET type = $1, credit = credit + 1, modified_at = UNIX_TIMESTAMP(NOW()) WHERE id > $2 AND name LIKE $3 AND (id_card IS NULL OR status IN ($4, $5, $6)) AND modified_at > created_at + $7 ORDER BY id ASC", result)
	assert.Equal(t, []interface{}{"sys", 1234, "%Du", 1, 2, 5, 86400}, args)
}
