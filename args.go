package pgsql

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Args stores arguments associated with a SQL.
type Args struct {
	namedArgs    map[string]int
	sqlNamedArgs map[string]int
	args         []interface{}
	onlyNamed    bool
}

func init() {
	// Predefine some $n args to avoid additional memory allocation.
	predefinedArgs = make([]string, 0, maxPredefinedArgs)

	for i := 0; i < maxPredefinedArgs; i++ {
		predefinedArgs = append(predefinedArgs, fmt.Sprintf("$%v", i))
	}
}

const maxPredefinedArgs = 64

var predefinedArgs []string

// Add adds an arg to Args and returns a placeholder.
func (args *Args) Add(arg interface{}) string {
	idx := args.add(arg)

	if idx < maxPredefinedArgs {
		return predefinedArgs[idx]
	}

	return fmt.Sprintf("$%v", idx)
}

func (args *Args) add(arg interface{}) int {
	idx := len(args.args)

	args.args = append(args.args, arg)
	return idx
}

// Compile compiles builder's format to standard sql and returns associated args.
//
// The format string uses a special syntax to represent arguments.
//
//	$? refers successive arguments passed in the call. It works similar as `%v` in `fmt.Sprintf`.
//	$0 $1 ... $n refers nth-argument passed in the call. Next $? will use arguments n+1.
//	${name} refers a named argument created by `Named` with `name`.
//	$$ is a "$" string.
func (args *Args) Compile(format string, initialValue ...interface{}) (query string, values []interface{}) {
	buf := &strings.Builder{}
	idx := strings.IndexRune(format, '$')
	offset := 0
	values = initialValue

	for idx >= 0 && len(format) > 0 {
		if idx > 0 {
			buf.WriteString(format[:idx])
		}

		format = format[idx+1:]

		// Treat the $ at the end of format is a normal $ rune.
		if len(format) == 0 {
			buf.WriteRune('$')
			break
		}

		if r := format[0]; r == '$' {
			buf.WriteRune('$')
			format = format[1:]
		} else if r == '{' {
			format, values = args.compileNamed(buf, format, values)
		} else if !args.onlyNamed && '0' <= r && r <= '9' {
			format, values, offset = args.compileDigits(buf, format, values, offset)
		} else if !args.onlyNamed && r == '?' {
			format, values, offset = args.compileSuccessive(buf, format[1:], values, offset)
		} else {
			// For unknown $ expression format, treat it as a normal $ rune.
			buf.WriteRune('$')
		}

		idx = strings.IndexRune(format, '$')
	}

	if len(format) > 0 {
		buf.WriteString(format)
	}

	query = buf.String()

	if len(args.sqlNamedArgs) > 0 {
		// Stabilize the sequence to make it easier to write test cases.
		ints := make([]int, 0, len(args.sqlNamedArgs))

		for _, p := range args.sqlNamedArgs {
			ints = append(ints, p)
		}

		sort.Ints(ints)

		for _, i := range ints {
			values = append(values, args.args[i])
		}
	}

	return
}

func (args *Args) compileNamed(buf *strings.Builder, format string, values []interface{}) (string, []interface{}) {
	i := 1

	for ; i < len(format) && format[i] != '}'; i++ {
		// Nothing.
	}

	// Invalid $ format. Ignore it.
	if i == len(format) {
		return format, values
	}

	name := format[1:i]
	format = format[i+1:]

	if p, ok := args.namedArgs[name]; ok {
		format, values, _ = args.compileSuccessive(buf, format, values, p)
	}

	return format, values
}

func (args *Args) compileDigits(buf *strings.Builder, format string, values []interface{}, offset int) (string, []interface{}, int) {
	i := 1

	for ; i < len(format) && '0' <= format[i] && format[i] <= '9'; i++ {
		// Nothing.
	}

	digits := format[:i]
	format = format[i:]

	if pointer, err := strconv.Atoi(digits); err == nil {
		return args.compileSuccessive(buf, format, values, pointer)
	}

	return format, values, offset
}

func (args *Args) compileSuccessive(buf *strings.Builder, format string, values []interface{}, offset int) (string, []interface{}, int) {
	if offset >= len(args.args) {
		return format, values, offset
	}

	arg := args.args[offset]
	values = args.compileArg(buf, values, arg)

	return format, values, offset + 1
}

func (args *Args) compileArg(buf *strings.Builder, values []interface{}, arg interface{}) []interface{} {
	switch a := arg.(type) {
	case Builder:
		var s string
		s, values = a.Build(values...)
		buf.WriteString(s)
	case rawArgs:
		buf.WriteString(a.expr)
	case listArgs:
		if len(a.args) > 0 {
			values = args.compileArg(buf, values, a.args[0])
		}

		for i := 1; i < len(a.args); i++ {
			buf.WriteString(", ")
			values = args.compileArg(buf, values, a.args[i])
		}
	default:
		fmt.Fprintf(buf, "$%d", len(values)+1)

		values = append(values, arg)
	}

	return values
}
