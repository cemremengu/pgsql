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
	args        *Args
	verb        string
	table       string
	returning   []string
	onConflict  []string
	assignments []string
	cols        []string
	values      [][]string
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

// Cols sets columns in INSERT.
func (ib *InsertBuilder) OnConflict(col ...string) *InsertBuilder {
	ib.onConflict = col
	return ib
}

// Cols sets columns in INSERT.
func (ib *InsertBuilder) DoUpdate(assignment ...string) *InsertBuilder {
	ib.assignments = assignment
	return ib
}

// Assign represents SET "field = value" in UPDATE.
func (ub *InsertBuilder) Set(col string) string {
	return fmt.Sprintf("%s = EXCLUDED.%s", col, col)
}

// Assign represents SET "field = value" in UPDATE.
func (ub *InsertBuilder) Assign(field string, value interface{}) string {
	return fmt.Sprintf("%s = %s", field, ub.args.Add(value))
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

	if len(ib.onConflict) != 0 {
		buf.WriteString(" ON CONFLICT (")
		buf.WriteString(strings.Join(ib.onConflict, ", "))
		buf.WriteString(")")

		if len(ib.assignments) != 0 {
			buf.WriteString(" DO UPDATE SET ")
			buf.WriteString(strings.Join(ib.assignments, ", "))
		}
	}

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
