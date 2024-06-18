Start server

```sh
go run examples/basic/main.go
```

Test server

```sh
curl -s http://localhost:3000/products/a32fb2bd-b402-4bea-93c2-4a0a567b2261 | jq

curl -s http://localhost:3000/products | jq
```