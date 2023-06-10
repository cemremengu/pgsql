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
	verb   string
	table  string
	cols   []string
	values [][]string

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

// // InsertIgnoreInto sets table name in INSERT IGNORE.
// func InsertIgnoreInto(table string) *InsertBuilder {
// 	return DefaultFlavor.NewInsertBuilder().InsertIgnoreInto(table)
// }

// // InsertIgnoreInto sets table name in INSERT IGNORE.
// func (ib *InsertBuilder) InsertIgnoreInto(table string) *InsertBuilder {
// 	ib.args.Flavor.PrepareInsertIgnore(table, ib)
// 	return ib
// }

// // ReplaceInto sets table name and changes the verb of ib to REPLACE.
// // REPLACE INTO is a MySQL extension to the SQL standard.
// func ReplaceInto(table string) *InsertBuilder {
// 	return DefaultFlavor.NewInsertBuilder().ReplaceInto(table)
// }

// ReplaceInto sets table name and changes the verb of ib to REPLACE.
// REPLACE INTO is a MySQL extension to the SQL standard.
// func (ib *InsertBuilder) ReplaceInto(table string) *InsertBuilder {
// 	ib.verb = "REPLACE"
// 	ib.table = Escape(table)
// 	ib.marker = insertMarkerAfterInsertInto
// 	return ib
// }

// Cols sets columns in INSERT.
func (ib *InsertBuilder) Cols(col ...string) *InsertBuilder {
	ib.cols = col
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

// String returns the compiled INSERT string.
func (ib *InsertBuilder) String() string {
	s, _ := ib.Build()
	return s
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

	return ib.args.Compile(buf.String(), initialArg...)
}

// Var returns a placeholder for value.
func (ib *InsertBuilder) Var(arg interface{}) string {
	return ib.args.Add(arg)
}
