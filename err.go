package errhandler

// HTTPError allows you to return custom status codes.
type HTTPError struct {
	status int
	err    error
}

// Error returns a new instance of Error.
func Error(status int, err error) HTTPError {
	return HTTPError{
		status: status,
		err:    err,
	}
}

func (err HTTPError) Error() string {
	return err.err.Error()
}
