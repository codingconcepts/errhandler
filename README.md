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
mux.Handle("GET /products", errhandler.Wrap(addProduct))
```

Update your existing handler signatures, by adding an `error` return type and utilizing the optional `errhandler.ParseJSON` and `errhandler.SendJSON` helper functions to reduce the amount of boilerplate code handlers require:

```go
func addProduct(w http.ResponseWriter, r *http.Request) error {
  var p product
  if err := errhandler.ParseJSON(r, &p); err != nil {
    return fmt.Errorf("parsing request json: %w", err)
  }

  products[p.ID] = p

  return errhandler.SendJSON(w, p)
}
```