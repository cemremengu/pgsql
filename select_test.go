package pgsql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkSelect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sb := NewSelectBuilder()
		sb.Select("*").
			From("demo.user").
			Where(sb.And(sb.Expr("test", "=", 1), sb.IsNotNull("deleted"))).
			Limit(10).
			Offset(1).
			Build()
	}
}

func TestSelect1(t *testing.T) {
	result, args := Select("id", "name").
		From("demo.user").
		Where("status = 1").
		Offset(1).
		Limit(10).
		Build()

	assert.Equal(t, "SELECT id, name FROM demo.user WHERE status = 1 LIMIT 10 OFFSET 1", result)
	assert.Empty(t, args)
}

func TestSelect2(t *testing.T) {
	result, args := Select("id", "name").
		From("demo.user").
		Where("status = 1").
		Build()

	assert.Equal(t, "SELECT id, name FROM demo.user WHERE status = 1", result)
	assert.Empty(t, args)
}

func TestSelect3(t *testing.T) {
	result, args := Select("*").
		From("demo.user").
		Build()

	assert.Equal(t, "SELECT * FROM demo.user", result)
	assert.Empty(t, args)
}

func TestSelect4(t *testing.T) {
	sb := NewSelectBuilder()

	result, args := sb.Select("*").
		From("demo.user").
		Where(sb.EQ("test", 1)).
		Build()

	assert.Equal(t, "SELECT * FROM demo.user WHERE test = $1", result)
	assert.Equal(t, []interface{}{1}, args)
}

func TestSelect5(t *testing.T) {
	sb := NewSelectBuilder()

	result, args := sb.Select("*").
		From("demo.user").
		Where(sb.Expr("test", ">", 1)).
		Build()

	assert.Equal(t, "SELECT * FROM demo.user WHERE test > $1", result)
	assert.Equal(t, []interface{}{1}, args)
}

func TestSelect6(t *testing.T) {
	sb := NewSelectBuilder()

	result, args := sb.Select("*").
		From("demo.user").
		Where(sb.And(sb.EQ("x", 2), sb.Expr("y", ">", 1))).
		Build()

	assert.Equal(t, "SELECT * FROM demo.user WHERE (x = $1 AND y > $2)", result)
	assert.Equal(t, []interface{}{2, 1}, args)
}

func TestSelect7(t *testing.T) {
	sb := NewSelectBuilder()

	result, args := sb.Select("*").
		From("demo.user").
		Where(sb.In("test", List([]int{1, 2}))).
		Build()

	assert.Equal(t, "SELECT * FROM demo.user WHERE test IN ($1, $2)", result)
	assert.Equal(t, []interface{}{1, 2}, args)
}

func TestSelect8(t *testing.T) {
	sb := NewSelectBuilder()

	result, args := sb.Select("*").
		From("demo.user").
		Join("demo.user_profile", "user.id = user_profile.user_id").
		Where(sb.In("test", List([]int{1, 2}))).
		Limit(10).
		Build()

	assert.Equal(t, "SELECT * FROM demo.user JOIN demo.user_profile ON user.id = user_profile.user_id WHERE test IN ($1, $2) LIMIT 10", result)
	assert.Equal(t, []interface{}{1, 2}, args)
}

func TestSelect9(t *testing.T) {
	sb := NewSelectBuilder()

	result, args := sb.Select("*").
		From("demo.user").
		Where(sb.And(sb.EQ("test", 1), sb.IsNotNull("deleted"))).
		Build()

	assert.Equal(t, "SELECT * FROM demo.user WHERE (test = $1 AND deleted IS NOT NULL)", result)
	assert.Equal(t, []interface{}{1}, args)
}
