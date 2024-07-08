package errsql

import (
	"fmt"
	"io"
	"runtime"
	"slices"
	"strings"

	"github.com/pkg/errors"
)

type stack []uintptr

var _ interface{ StackTrace() errors.StackTrace } = (*stack)(nil)

func (s *stack) Format(st fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case st.Flag('+'):
			for _, pc := range s.external() {
				f := errors.Frame(pc)
				_, _ = fmt.Fprintf(st, "\n%+v", f)
			}
		}
	}
}

func (s *stack) StackTrace() errors.StackTrace {
	stack := s.external()
	f := make([]errors.Frame, len(stack))
	for i := 0; i < len(f); i++ {
		f[i] = errors.Frame((stack)[i])
	}
	return f
}

// external returns only the stack frames external to this library, and
// the stdlib.
func (s *stack) external() []uintptr {
	// We must find at least _one_ internal frame, to account for the error
	// handler function itself.
	var foundInternal bool
	if start := slices.IndexFunc(*s, func(pc uintptr) bool {
		if isExternal(pc) {
			return foundInternal
		}
		foundInternal = true
		return false
	}); start > 0 {
		return (*s)[start:]
	}
	return *s
}

// isExternal returns true if pc represents a function outside of this and the `database/sql` packages.
func isExternal(pc uintptr) bool {
	name := runtime.FuncForPC(pc - 1).Name()
	switch {
	case strings.HasPrefix(name, "gitlab.com/flimzy/errsql."),
		strings.HasPrefix(name, "database/sql."):
		return false
	default:
		return true
	}
}

// AddStacktrace adds a [github.com/pkg/error.Stacktrace] to the error at the
// point at which the db function is called. It excludes stack frames in the
// `database/sql` package, and in this package, as well as from the error
// handler function.
func AddStacktrace(err error) error {
	if err == nil {
		return nil
	}
	return &withStack{
		err,
		callers(),
	}
}

func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[0:n]
	return &st
}

type withStack struct {
	error
	*stack
}

func (w *withStack) Unwrap() error { return w.error }

func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v", w.Unwrap())
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, w.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", w.Error())
	}
}
