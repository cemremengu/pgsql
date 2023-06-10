package pgsql

// Builder is a general SQL builder.
// It's used by Args to create nested SQL like the `IN` expression in
// `SELECT * FROM t1 WHERE id IN (SELECT id FROM t2)`.
type Builder interface {
	Build(initialArg ...interface{}) (sql string, args []interface{})
}

type builder struct {
	args   *Args
	format string
}

func (cb *builder) Build(initialArg ...interface{}) (sql string, args []interface{}) {
	return cb.args.Compile(cb.format, initialArg...)
}

// Build creates a Builder from a format string.
// The format string uses special syntax to represent arguments.
// See doc in `Args#Compile` for syntax details.
func Build(format string, arg ...interface{}) Builder {
	args := &Args{}

	for _, a := range arg {
		args.Add(a)
	}

	return &builder{
		args:   args,
		format: format,
	}
}
