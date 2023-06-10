package pgsql

import (
	"strconv"
	"strings"
)

const (
	unionDistinct = " UNION " // Default union type is DISTINCT.
	unionAll      = " UNION ALL "
)

// NewUnionBuilder creates a new UNION builder.
func NewUnionBuilder() *UnionBuilder {
	return newUnionBuilder()
}

func newUnionBuilder() *UnionBuilder {
	return &UnionBuilder{
		limit:  -1,
		offset: -1,

		args: &Args{},
	}
}

// UnionBuilder is a builder to build UNION.
type UnionBuilder struct {
	opt         string
	builders    []Builder
	orderByCols []string
	order       string
	limit       int
	offset      int

	args *Args
}

// Union unions all builders together using UNION operator.
func Union(builders ...Builder) *UnionBuilder {
	return NewUnionBuilder().Union(builders...)
}

// Union unions all builders together using UNION operator.
func (ub *UnionBuilder) Union(builders ...Builder) *UnionBuilder {
	return ub.union(unionDistinct, builders...)
}

// UnionAll unions all builders together using UNION ALL operator.
func UnionAll(builders ...Builder) *UnionBuilder {
	return NewUnionBuilder().UnionAll(builders...)
}

// UnionAll unions all builders together using UNION ALL operator.
func (ub *UnionBuilder) UnionAll(builders ...Builder) *UnionBuilder {
	return ub.union(unionAll, builders...)
}

func (ub *UnionBuilder) union(opt string, builders ...Builder) *UnionBuilder {
	ub.opt = opt
	ub.builders = builders
	return ub
}

// OrderBy sets columns of ORDER BY in SELECT.
func (ub *UnionBuilder) OrderBy(col ...string) *UnionBuilder {
	ub.orderByCols = col
	return ub
}

// Asc sets order of ORDER BY to ASC.
func (ub *UnionBuilder) Asc() *UnionBuilder {
	ub.order = "ASC"
	return ub
}

// Desc sets order of ORDER BY to DESC.
func (ub *UnionBuilder) Desc() *UnionBuilder {
	ub.order = "DESC"
	return ub
}

// Limit sets the LIMIT in SELECT.
func (ub *UnionBuilder) Limit(limit int) *UnionBuilder {
	ub.limit = limit
	return ub
}

// Offset sets the LIMIT offset in SELECT.
func (ub *UnionBuilder) Offset(offset int) *UnionBuilder {
	ub.offset = offset
	return ub
}

// Build returns compiled SELECT string and args.
// They can be used in `DB#Query` of package `database/sql` directly.
func (ub *UnionBuilder) Build(initialArg ...interface{}) (sql string, args []interface{}) {
	buf := &strings.Builder{}

	if len(ub.builders) > 0 {

		buf.WriteRune('(')

		buf.WriteString(ub.Var(ub.builders[0]))

		buf.WriteRune(')')

		for _, b := range ub.builders[1:] {
			buf.WriteString(ub.opt)

			buf.WriteRune('(')

			buf.WriteString(ub.Var(b))

			buf.WriteRune(')')
		}
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

	if ub.offset >= 0 {
		buf.WriteString(" OFFSET ")
		buf.WriteString(strconv.Itoa(ub.offset))
	}

	return ub.args.Compile(buf.String(), initialArg...)
}

// Var returns a placeholder for value.
func (ub *UnionBuilder) Var(arg interface{}) string {
	return ub.args.Add(arg)
}
