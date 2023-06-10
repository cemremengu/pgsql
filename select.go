package sqlb

import (
	"strconv"
	"strings"
)

// SelectBuilder is a builder to build SELECT.
type SelectBuilder struct {
	// Cond

	distinct   bool
	tables     []string
	selectCols []string
	// joinOptions []JoinOption
	joinTables  []string
	joinExprs   [][]string
	whereExprs  []string
	havingExprs []string
	groupByCols []string
	orderByCols []string
	order       string
	limit       int
	offset      int
	forWhat     string

	args *Args
}

func newSelectBuilder() *SelectBuilder {
	args := &Args{}
	return &SelectBuilder{
		// Cond: Cond{
		// 	Args: args,
		// },
		limit:  -1,
		offset: -1,
		args:   args,
	}
}

func Select(col ...string) *SelectBuilder {
	return newSelectBuilder().Select(col...)
}

func (sb *SelectBuilder) Select(col ...string) *SelectBuilder {
	sb.selectCols = col
	return sb
}

// From sets table names in SELECT.
func (sb *SelectBuilder) From(table ...string) *SelectBuilder {
	sb.tables = table
	return sb
}

// BuildWithFlavor returns compiled SELECT string and args with flavor and initial args.
// They can be used in `DB#Query` of package `database/sql` directly.
func (sb *SelectBuilder) Build(initialArg ...interface{}) (sql string, args []interface{}) {
	buf := &strings.Builder{}
	buf.WriteString("SELECT ")

	if sb.distinct {
		buf.WriteString("DISTINCT ")
	}

	buf.WriteString(strings.Join(sb.selectCols, ", "))

	buf.WriteString(" FROM ")
	buf.WriteString(strings.Join(sb.tables, ", "))

	// for i := range sb.joinTables {
	// 	if option := sb.joinOptions[i]; option != "" {
	// 		buf.WriteRune(' ')
	// 		buf.WriteString(string(option))
	// 	}

	// 	buf.WriteString(" JOIN ")
	// 	buf.WriteString(sb.joinTables[i])

	// 	if exprs := sb.joinExprs[i]; len(exprs) > 0 {
	// 		buf.WriteString(" ON ")
	// 		buf.WriteString(strings.Join(sb.joinExprs[i], " AND "))
	// 	}

	// }

	if len(sb.whereExprs) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(sb.whereExprs, " AND "))

	}

	if len(sb.groupByCols) > 0 {
		buf.WriteString(" GROUP BY ")
		buf.WriteString(strings.Join(sb.groupByCols, ", "))

		if len(sb.havingExprs) > 0 {
			buf.WriteString(" HAVING ")
			buf.WriteString(strings.Join(sb.havingExprs, " AND "))
		}

	}

	if len(sb.orderByCols) > 0 {
		buf.WriteString(" ORDER BY ")
		buf.WriteString(strings.Join(sb.orderByCols, ", "))

		if sb.order != "" {
			buf.WriteRune(' ')
			buf.WriteString(sb.order)
		}

	}

	if sb.limit >= 0 {
		buf.WriteString(" LIMIT ")
		buf.WriteString(strconv.Itoa(sb.limit))
	}

	if sb.offset >= 0 {
		buf.WriteString(" OFFSET ")
		buf.WriteString(strconv.Itoa(sb.offset))
	}

	if sb.forWhat != "" {
		buf.WriteString(" FOR ")
		buf.WriteString(sb.forWhat)

	}

	return sb.args.CompileWithFlavor(buf.String(), initialArg...)
}
