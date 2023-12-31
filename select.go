package pgsql

import (
	"fmt"
	"strconv"
	"strings"
)

// JoinOption is the option in JOIN.
type JoinOption string

// Join options.
const (
	FullJoin       JoinOption = "FULL"
	FullOuterJoin  JoinOption = "FULL OUTER"
	InnerJoin      JoinOption = "INNER"
	LeftJoin       JoinOption = "LEFT"
	LeftOuterJoin  JoinOption = "LEFT OUTER"
	RightJoin      JoinOption = "RIGHT"
	RightOuterJoin JoinOption = "RIGHT OUTER"
)

// SelectBuilder is a builder to build SELECT.
type SelectBuilder struct {
	Cond

	args    *Args
	forWhat string

	order       string
	havingExprs []string
	joinTables  []string
	joinExprs   [][]string
	whereExprs  []string
	joinOptions []JoinOption
	groupByCols []string
	orderByCols []string
	selectCols  []string
	tables      []string
	limit       int
	offset      int

	distinct bool
}

// NewSelectBuilder creates a new SELECT builder.
func NewSelectBuilder() *SelectBuilder {
	return newSelectBuilder()
}

func newSelectBuilder() *SelectBuilder {
	args := &Args{}
	return &SelectBuilder{
		Cond: Cond{
			Args: args,
		},
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

func (sb *SelectBuilder) Where(andExpr ...string) *SelectBuilder {
	sb.whereExprs = append(sb.whereExprs, andExpr...)
	return sb
}

func (sb *SelectBuilder) Limit(limit int) *SelectBuilder {
	sb.limit = limit
	return sb
}

// Offset sets the LIMIT offset in SELECT.
func (sb *SelectBuilder) Offset(offset int) *SelectBuilder {
	sb.offset = offset
	return sb
}

// Having sets expressions of HAVING in SELECT.
func (sb *SelectBuilder) Having(andExpr ...string) *SelectBuilder {
	sb.havingExprs = append(sb.havingExprs, andExpr...)
	return sb
}

// GroupBy sets columns of GROUP BY in SELECT.
func (sb *SelectBuilder) GroupBy(col ...string) *SelectBuilder {
	sb.groupByCols = append(sb.groupByCols, col...)
	return sb
}

// OrderBy sets columns of ORDER BY in SELECT with the provided order.
func (sb *SelectBuilder) OrderBy(order string, col ...string) *SelectBuilder {
	sb.order = order
	sb.orderByCols = append(sb.orderByCols, col...)
	return sb
}

// OrderByAsc sets columns of ORDER BY ASC in SELECT.
func (sb *SelectBuilder) OrderByAsc(col ...string) *SelectBuilder {
	sb.order = "ASC"
	sb.orderByCols = append(sb.orderByCols, col...)
	return sb
}

// OrderByDesc sets columns of ORDER BY DESC in SELECT.
func (sb *SelectBuilder) OrderByDesc(col ...string) *SelectBuilder {
	sb.order = "DESC"
	sb.orderByCols = append(sb.orderByCols, col...)
	return sb
}

// BuilderAs returns an AS expression wrapping a complex SQL.
// According to SQL syntax, SQL built by builder is surrounded by parens.
func (sb *SelectBuilder) BuilderAs(builder Builder, alias string) string {
	return fmt.Sprintf("(%s) AS %s", sb.Var(builder), alias)
}

// // Asc sets order of ORDER BY to ASC.
// func (sb *SelectBuilder) Asc() *SelectBuilder {
// 	sb.order = "ASC"
// 	return sb
// }

// // Desc sets order of ORDER BY to DESC.
// func (sb *SelectBuilder) Desc() *SelectBuilder {
// 	sb.order = "DESC"
// 	return sb
// }

// Join sets expressions of JOIN in SELECT.
//
// It builds a JOIN expression like
//
//	JOIN table ON onExpr[0] AND onExpr[1] ...
func (sb *SelectBuilder) Join(table string, onExpr ...string) *SelectBuilder {
	return sb.JoinWithOption("", table, onExpr...)
}

func (sb *SelectBuilder) LeftJoin(table string, onExpr ...string) *SelectBuilder {
	return sb.JoinWithOption(LeftJoin, table, onExpr...)
}

// JoinWithOption sets expressions of JOIN with an option.
//
// It builds a JOIN expression like
//
//	option JOIN table ON onExpr[0] AND onExpr[1] ...
//
// Here is a list of supported options.
//   - FullJoin: FULL JOIN
//   - FullOuterJoin: FULL OUTER JOIN
//   - InnerJoin: INNER JOIN
//   - LeftJoin: LEFT JOIN
//   - LeftOuterJoin: LEFT OUTER JOIN
//   - RightJoin: RIGHT JOIN
//   - RightOuterJoin: RIGHT OUTER JOIN
func (sb *SelectBuilder) JoinWithOption(option JoinOption, table string, onExpr ...string) *SelectBuilder {
	sb.joinOptions = append(sb.joinOptions, option)
	sb.joinTables = append(sb.joinTables, table)
	sb.joinExprs = append(sb.joinExprs, onExpr)
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

	for i := range sb.joinTables {
		if option := sb.joinOptions[i]; option != "" {
			buf.WriteRune(' ')
			buf.WriteString(string(option))
		}

		buf.WriteString(" JOIN ")
		buf.WriteString(sb.joinTables[i])

		if exprs := sb.joinExprs[i]; len(exprs) > 0 {
			buf.WriteString(" ON ")
			buf.WriteString(strings.Join(sb.joinExprs[i], " AND "))
		}

	}

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

	return sb.args.Compile(buf.String(), initialArg...)
}
