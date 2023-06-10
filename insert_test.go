package pgsql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsert1(t *testing.T) {
	result, args := InsertInto("demo.user").
		Cols("id", "name", "status").
		Values(4, "Sample", 2).
		Build()

	assert.Equal(t, "INSERT INTO demo.user (id, name, status) VALUES ($1, $2, $3)", result)
	assert.Equal(t, []interface{}{4, "Sample", 2}, args)
}

func TestInsert2(t *testing.T) {
	ib := NewInsertBuilder()
	ib.InsertInto("demo.user")
	ib.Cols("id", "name", "status", "created_at", "updated_at")
	ib.Values(1, "Huan Du", 1, Raw("UNIX_TIMESTAMP(NOW())"))
	ib.Values(2, "Charmy Liu", 1, 1234567890)

	result, args := ib.Build()
	assert.Equal(t, "INSERT INTO demo.user (id, name, status, created_at, updated_at) VALUES ($1, $2, $3, UNIX_TIMESTAMP(NOW())), ($4, $5, $6, $7)", result)
	assert.Equal(t, []interface{}{1, "Huan Du", 1, 2, "Charmy Liu", 1, 1234567890}, args)
}
