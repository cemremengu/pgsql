package sqlb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
