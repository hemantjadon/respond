# respond

Package responds provides low touch, minimal API for sending HTTP API responses
in go.

For simple string responses can be used simply as:
```go
func handler(w http.ResponseWriter, r *http.Request) {
     respond.With(w, http.StatusOK, []byte(`Hello World!`))
}
```

For more complex use cases where we want to send JSON across this respond
provides handy utility function which can be used as follows:
```go
type response struct {
    Message string `json: "message"`
}

func handler(w http.ResponseWriter, r *http.Request) {
    resp := response{Message: "Hello World!"}
    respond.WithJSON(w, http.StatusOK, response)
}
```
While sending JSON responses correct HTTP `Content-Type: applocation/json; utf-8`
is also set.
