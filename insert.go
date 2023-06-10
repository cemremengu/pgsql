package pgsql

import (
	"fmt"
	"strings"
)

// NewInsertBuilder creates a new INSERT builder.
func NewInsertBuilder() *InsertBuilder {
	return newInsertBuilder()
}

func newInsertBuilder() *InsertBuilder {
	args := &Args{}
	return &InsertBuilder{
		verb: "INSERT",
		args: args,
	}
}

// InsertBuilder is a builder to build INSERT.
type InsertBuilder struct {
	verb      string
	returning []string
	table     string
	cols      []string
	values    [][]string

	args *Args
}

// InsertInto sets table name in INSERT.
func InsertInto(table string) *InsertBuilder {
	return NewInsertBuilder().InsertInto(table)
}

// InsertInto sets table name in INSERT.
func (ib *InsertBuilder) InsertInto(table string) *InsertBuilder {
	ib.table = table
	return ib
}

// Cols sets columns in INSERT.
func (ib *InsertBuilder) Cols(col ...string) *InsertBuilder {
	ib.cols = col
	return ib
}

func (ib *InsertBuilder) Returning(col ...string) *InsertBuilder {
	ib.returning = col
	return ib
}

// Values adds a list of values for a row in INSERT.
func (ib *InsertBuilder) Values(value ...interface{}) *InsertBuilder {
	placeholders := make([]string, 0, len(value))

	for _, v := range value {
		placeholders = append(placeholders, ib.args.Add(v))
	}

	ib.values = append(ib.values, placeholders)
	return ib
}

// BuildWithFlavor returns compiled INSERT string and args with flavor and initial args.
// They can be used in `DB#Query` of package `database/sql` directly.
func (ib *InsertBuilder) Build(initialArg ...interface{}) (sql string, args []interface{}) {
	buf := &strings.Builder{}
	buf.WriteString(ib.verb)
	buf.WriteString(" INTO ")
	buf.WriteString(ib.table)

	if len(ib.cols) > 0 {
		buf.WriteString(" (")
		buf.WriteString(strings.Join(ib.cols, ", "))
		buf.WriteString(")")

	}

	buf.WriteString(" VALUES ")
	values := make([]string, 0, len(ib.values))

	for _, v := range ib.values {
		values = append(values, fmt.Sprintf("(%v)", strings.Join(v, ", ")))
	}

	buf.WriteString(strings.Join(values, ", "))

	if len(ib.returning) != 0 {
		buf.WriteString(" RETURNING ")
		buf.WriteString(strings.Join(ib.returning, ", "))
	}

	return ib.args.Compile(buf.String(), initialArg...)
}

// Var returns a placeholder for value.
func (ib *InsertBuilder) Var(arg interface{}) string {
	return ib.args.Add(arg)
}
