package sqlb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelect1(t *testing.T) {
	result, args := Select("col").From("table").Build()

	assert.Equal(t, "SELECT col FROM table", result)
	assert.Empty(t, args)

	fmt.Println(result)
}

func TestSelect2(t *testing.T) {
	result, args := Select("col").From("table").Build()

	assert.Equal(t, "SELECT col FROM table", result)
	assert.Empty(t, args)

	fmt.Println(result)
}
