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

func TestInsert3(t *testing.T) {
	ib := NewInsertBuilder()
	ib.InsertInto("demo.user")
	ib.Cols("id", "name", "status", "created_at", "updated_at")
	ib.Values(1, "Charmy Liu", 1, 1234567890)
	ib.Returning("id", "name")

	result, args := ib.Build()
	assert.Equal(t, "INSERT INTO demo.user (id, name, status, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id, name", result)
	assert.Equal(t, []interface{}{1, "Charmy Liu", 1, 1234567890}, args)
}

func TestInsert4(t *testing.T) {
	ib := NewInsertBuilder()
	ib.InsertInto("demo.user")
	ib.Cols("id", "name", "status", "created_at", "updated_at")
	ib.Values(1, "Charmy Liu", 1, 1234567890)
	ib.OnConflict("id", "name")
	ib.DoUpdate(ib.Assign("status", 2), ib.Assign("updated_at", 12345))

	result, args := ib.Build()
	assert.Equal(t, "INSERT INTO demo.user (id, name, status, created_at, updated_at) VALUES ($1, $2, $3, $4) ON CONFLICT (id, name) DO UPDATE SET status = $5, updated_at = $6", result)
	assert.Equal(t, []interface{}{1, "Charmy Liu", 1, 1234567890, 2, 12345}, args)
}

func TestInsert5(t *testing.T) {
	ib := NewInsertBuilder()
	ib.InsertInto("demo.user")
	ib.Cols("id", "name", "status", "created_at", "updated_at")
	ib.Values(1, "Charmy Liu", 1, 1234567890)
	ib.OnConflict("id", "name")
	ib.DoUpdate(ib.Set("status"), ib.Set("updated_at"))

	result, args := ib.Build()
	assert.Equal(t, "INSERT INTO demo.user (id, name, status, created_at, updated_at) VALUES ($1, $2, $3, $4) ON CONFLICT (id, name) DO UPDATE SET status = EXCLUDED.status, updated_at = EXCLUDED.updated_at", result)
	assert.Equal(t, []interface{}{1, "Charmy Liu", 1, 1234567890}, args)
}
