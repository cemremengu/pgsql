package pgsql

// Builder is a general SQL builder.
// It's used by Args to create nested SQL like the `IN` expression in
// `SELECT * FROM t1 WHERE id IN (SELECT id FROM t2)`.
type Builder interface {
	Build() (sql string, args []interface{})
	BuildWithFlavor(initialArg ...interface{}) (sql string, args []interface{})
}

type compiledBuilder struct {
	args   *Args
	format string
}

var _ Builder = new(compiledBuilder)

func (cb *compiledBuilder) Build() (sql string, args []interface{}) {
	return cb.args.Compile(cb.format)
}

func (cb *compiledBuilder) BuildWithFlavor(initialArg ...interface{}) (sql string, args []interface{}) {
	return cb.args.Compile(cb.format, initialArg...)
}

type flavoredBuilder struct {
	builder Builder
}

func (fb *flavoredBuilder) Build() (sql string, args []interface{}) {
	return fb.builder.BuildWithFlavor()
}

func (fb *flavoredBuilder) BuildWithFlavor(initialArg ...interface{}) (sql string, args []interface{}) {
	return fb.builder.BuildWithFlavor(initialArg...)
}

// WithFlavor creates a new Builder based on builder with a default flavor.
func WithFlavor(builder Builder) Builder {
	return &flavoredBuilder{
		builder: builder,
	}
}

// Build creates a Builder from a format string.
// The format string uses special syntax to represent arguments.
// See doc in `Args#Compile` for syntax details.
func Build(format string, arg ...interface{}) Builder {
	args := &Args{}

	for _, a := range arg {
		args.Add(a)
	}

	return &compiledBuilder{
		args:   args,
		format: format,
	}
}
