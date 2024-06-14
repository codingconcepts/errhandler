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
mux.Handle("GET /products/{id}", errhandler.Wrap(getProduct))
```

Update your existing handler signature, by adding an `error` return type:

```go
func addProduct(w http.ResponseWriter, r *http.Request) error {
	...
}
```

Utilize the `ParseJSON` and `SendJSON` functions to reduce the amount of boilerplate code handlers require:

``` go
var p product
if err := errhandler.ParseJSON(r, &p); err != nil {
  return fmt.Errorf("parsing request json: %w", err)
}

products[p.ID] = p

return errhandler.SendJSON(w, p)
```