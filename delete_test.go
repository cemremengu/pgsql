package pgsql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelete1(t *testing.T) {
	result, args := DeleteFrom("demo.user").
		Where(
			"status = 1",
		).
		Limit(10).
		Build()

	assert.Equal(t, "DELETE FROM demo.user WHERE status = 1 LIMIT 10", result)
	assert.Empty(t, args)
}

func TestDelete2(t *testing.T) {
	db := NewDeleteBuilder()
	db.DeleteFrom("demo.user")
	db.Where(
		db.GT("id", 1234),
		db.Like("name", "%Du"),
		db.Or(
			db.IsNull("id_card"),
			db.In("status", 1, 2, 5),
		),
		"modified_at > created_at + "+db.Var(86400), // It's allowed to write arbitrary SQL.
	)

	result, args := db.Build()

	assert.Equal(t, "DELETE FROM demo.user WHERE id > $1 AND name LIKE $2 AND (id_card IS NULL OR status IN ($3, $4, $5)) AND modified_at > created_at + $6", result)
	assert.Equal(t, []interface{}{1234, "%Du", 1, 2, 5, 86400}, args)
}
