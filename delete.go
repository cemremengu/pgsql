package pgsql

import (
	"strconv"
	"strings"
)

// NewDeleteBuilder creates a new DELETE builder.
func NewDeleteBuilder() *DeleteBuilder {
	return newDeleteBuilder()
}

func newDeleteBuilder() *DeleteBuilder {
	args := &Args{}
	return &DeleteBuilder{
		Cond: Cond{
			Args: args,
		},
		limit: -1,
		args:  args,
	}
}

// DeleteBuilder is a builder to build DELETE.
type DeleteBuilder struct {
	Cond

	args *Args

	table       string
	order       string
	returning   []string
	whereExprs  []string
	orderByCols []string
	limit       int
}

// DeleteFrom sets table name in DELETE.
func DeleteFrom(table string) *DeleteBuilder {
	return NewDeleteBuilder().DeleteFrom(table)
}

// DeleteFrom sets table name in DELETE.
func (db *DeleteBuilder) DeleteFrom(table string) *DeleteBuilder {
	db.table = table
	return db
}

// Where sets expressions of WHERE in DELETE.
func (db *DeleteBuilder) Where(andExpr ...string) *DeleteBuilder {
	db.whereExprs = append(db.whereExprs, andExpr...)
	return db
}

// OrderBy sets columns of ORDER BY in DELETE.
func (db *DeleteBuilder) OrderBy(col ...string) *DeleteBuilder {
	db.orderByCols = col
	return db
}

// Asc sets order of ORDER BY to ASC.
func (db *DeleteBuilder) Asc() *DeleteBuilder {
	db.order = "ASC"
	return db
}

// Desc sets order of ORDER BY to DESC.
func (db *DeleteBuilder) Desc() *DeleteBuilder {
	db.order = "DESC"
	return db
}

// Limit sets the LIMIT in DELETE.
func (db *DeleteBuilder) Limit(limit int) *DeleteBuilder {
	db.limit = limit
	return db
}

// String returns the compiled DELETE string.
func (db *DeleteBuilder) String() string {
	s, _ := db.Build()
	return s
}

// Build returns compiled DELETE string and args.
// They can be used in `DB#Query` of package `database/sql` directly.
func (db *DeleteBuilder) Build(initialArg ...interface{}) (sql string, args []interface{}) {
	buf := &strings.Builder{}
	buf.WriteString("DELETE FROM ")
	buf.WriteString(db.table)

	if len(db.whereExprs) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(db.whereExprs, " AND "))

	}

	if len(db.orderByCols) > 0 {
		buf.WriteString(" ORDER BY ")
		buf.WriteString(strings.Join(db.orderByCols, ", "))

		if db.order != "" {
			buf.WriteRune(' ')
			buf.WriteString(db.order)
		}

	}

	if db.limit >= 0 {
		buf.WriteString(" LIMIT ")
		buf.WriteString(strconv.Itoa(db.limit))
	}

	if len(db.returning) != 0 {
		buf.WriteString(" RETURNING ")
		buf.WriteString(strings.Join(db.returning, ", "))
	}

	return db.args.Compile(buf.String(), initialArg...)
}
