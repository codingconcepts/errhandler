package errhandler

// Middleware is a helper for adding middleware to wrapped handelers.
type Middleware func(Wrap) Wrap

// Chain multiple middleware functions together into a single middleware.
func Chain(m ...Middleware) Middleware {
	return func(n Wrap) Wrap {
		for i := len(m) - 1; i >= 0; i-- {
			n = m[i](n)
		}

		return n
	}
}
