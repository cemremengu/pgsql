package pgsql

import (
	"fmt"
	"strconv"
	"strings"
)

// NewUpdateBuilder creates a new UPDATE builder.
func NewUpdateBuilder() *UpdateBuilder {
	return newUpdateBuilder()
}

func newUpdateBuilder() *UpdateBuilder {
	args := &Args{}
	return &UpdateBuilder{
		Cond: Cond{
			Args: args,
		},
		limit: -1,
		args:  args,
	}
}

// UpdateBuilder is a builder to build UPDATE.
type UpdateBuilder struct {
	Cond

	table       string
	assignments []string
	whereExprs  []string
	orderByCols []string
	order       string
	limit       int

	args *Args
}

// Update sets table name in UPDATE.
func Update(table string) *UpdateBuilder {
	return NewUpdateBuilder().Update(table)
}

// Update sets table name in UPDATE.
func (ub *UpdateBuilder) Update(table string) *UpdateBuilder {
	ub.table = table
	return ub
}

// Set sets the assignments in SET.
func (ub *UpdateBuilder) Set(assignment ...string) *UpdateBuilder {
	ub.assignments = assignment
	return ub
}

// SetMore appends the assignments in SET.
func (ub *UpdateBuilder) SetMore(assignment ...string) *UpdateBuilder {
	ub.assignments = append(ub.assignments, assignment...)
	return ub
}

// Where sets expressions of WHERE in UPDATE.
func (ub *UpdateBuilder) Where(andExpr ...string) *UpdateBuilder {
	ub.whereExprs = append(ub.whereExprs, andExpr...)
	return ub
}

// Assign represents SET "field = value" in UPDATE.
func (ub *UpdateBuilder) Assign(field string, value interface{}) string {
	return fmt.Sprintf("%s = %s", field, ub.args.Add(value))
}

// Incr represents SET "field = field + 1" in UPDATE.
func (ub *UpdateBuilder) Incr(field string) string {
	f := field
	return fmt.Sprintf("%s = %s + 1", f, f)
}

// Decr represents SET "field = field - 1" in UPDATE.
func (ub *UpdateBuilder) Decr(field string) string {
	f := field
	return fmt.Sprintf("%s = %s - 1", f, f)
}

// Add represents SET "field = field + value" in UPDATE.
func (ub *UpdateBuilder) Add(field string, value interface{}) string {
	f := field
	return fmt.Sprintf("%s = %s + %s", f, f, ub.args.Add(value))
}

// Sub represents SET "field = field - value" in UPDATE.
func (ub *UpdateBuilder) Sub(field string, value interface{}) string {
	f := field
	return fmt.Sprintf("%s = %s - %s", f, f, ub.args.Add(value))
}

// Mul represents SET "field = field * value" in UPDATE.
func (ub *UpdateBuilder) Mul(field string, value interface{}) string {
	f := field
	return fmt.Sprintf("%s = %s * %s", f, f, ub.args.Add(value))
}

// Div represents SET "field = field / value" in UPDATE.
func (ub *UpdateBuilder) Div(field string, value interface{}) string {
	f := field
	return fmt.Sprintf("%s = %s / %s", f, f, ub.args.Add(value))
}

// OrderBy sets columns of ORDER BY in UPDATE.
func (ub *UpdateBuilder) OrderBy(col ...string) *UpdateBuilder {
	ub.orderByCols = col
	return ub
}

// Asc sets order of ORDER BY to ASC.
func (ub *UpdateBuilder) Asc() *UpdateBuilder {
	ub.order = "ASC"
	return ub
}

// Desc sets order of ORDER BY to DESC.
func (ub *UpdateBuilder) Desc() *UpdateBuilder {
	ub.order = "DESC"
	return ub
}

// Limit sets the LIMIT in UPDATE.
func (ub *UpdateBuilder) Limit(limit int) *UpdateBuilder {
	ub.limit = limit
	return ub
}

// String returns the compiled UPDATE string.
func (ub *UpdateBuilder) String() string {
	s, _ := ub.Build()
	return s
}

// BuildWithFlavor returns compiled UPDATE string and args with flavor and initial args.
// They can be used in `DB#Query` of package `database/sql` directly.
func (ub *UpdateBuilder) Build(initialArg ...interface{}) (sql string, args []interface{}) {
	buf := &strings.Builder{}
	buf.WriteString("UPDATE ")
	buf.WriteString(ub.table)

	buf.WriteString(" SET ")
	buf.WriteString(strings.Join(ub.assignments, ", "))

	if len(ub.whereExprs) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(ub.whereExprs, " AND "))
	}

	if len(ub.orderByCols) > 0 {
		buf.WriteString(" ORDER BY ")
		buf.WriteString(strings.Join(ub.orderByCols, ", "))

		if ub.order != "" {
			buf.WriteRune(' ')
			buf.WriteString(ub.order)
		}

	}

	if ub.limit >= 0 {
		buf.WriteString(" LIMIT ")
		buf.WriteString(strconv.Itoa(ub.limit))

	}

	return ub.args.Compile(buf.String(), initialArg...)
}
