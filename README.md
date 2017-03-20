# simple-api

Small bootstrap project which counts and records api stats and rpm.

## Run

Start the server:

```go run api.go views.go middleware.go cache.go```

and check the rpm with

```curl http://localhost:8080/```


## Tests

```go test . -v```
