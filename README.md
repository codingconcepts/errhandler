# errhandler
Wraps Go's http.HandlerFunc, allowing for simpler use of http.NewServeMux().

### Install

```sh
go get github.com/codingconcepts/errhandler
```

### Usage

Wrap Go's stdlib `http.HandlerFunc` using `errhandler.Wrap()`:

```go
mux := http.NewServeMux()
mux.Handle("POST /products", errhandler.Wrap(addProduct))
```

Update your existing handler signatures, by adding an `error` return type and utilizing the optional `errhandler.ParseJSON` and `errhandler.SendJSON` helper functions to reduce the amount of boilerplate code handlers require:

```go
func addProduct(w http.ResponseWriter, r *http.Request) error {
  var p product
  if err := errhandler.ParseJSON(r, &p); err != nil {
    return err

    // Or, if you'd prefer to customise the status code:
    // return errhandler.Error(http.StatusUnprocessableEntity, err)
  }
  
  products[p.ID] = p
  
  return errhandler.SendJSON(w, p)
}
```

### Middleware

errhandler contains helper objects for building middleware

* A `errhandler.Middleware` type, which is simply a function that takes a Wrap function and returns a Wrap function:

```go
mux := http.NewServeMux()
mux.Handle("GET /products/{id}", errhandler.Wrap(midLog(getProduct)))

...

func midLog(n errhandler.Wrap) errhandler.Wrap {
	return func(w http.ResponseWriter, r *http.Request) error {
		log.Printf("1 %s %s", r.Method, r.URL.Path)
		return n(w, r)
	}
}

func addProduct(w http.ResponseWriter, r *http.Request) error {
  ...
}
```

* A `errhandler.Chain` function, which allows `Middleware` functions easy to chain:

```go
chain := errhandler.Chain(midLog1, midLog2)

mux := http.NewServeMux()
mux.Handle("GET /products/{id}", errhandler.Wrap(chain(getProducts)))

func midLog1(n errhandler.Wrap) errhandler.Wrap {
	return func(w http.ResponseWriter, r *http.Request) error {
		log.Printf("1 %s %s", r.Method, r.URL.Path)
		return n(w, r)
	}
}

func midLog2(n errhandler.Wrap) errhandler.Wrap {
	return func(w http.ResponseWriter, r *http.Request) error {
		log.Printf("2 %s %s", r.Method, r.URL.Path)
		return n(w, r)
	}
}

func addProduct(w http.ResponseWriter, r *http.Request) error {
  ...
}
```