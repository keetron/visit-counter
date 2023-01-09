# Visitor counter

This microservice allows for a POST call to be made to the `/visit` endpoint to add to the url visit counter, counting each visitor only once.

To start the service, use `go run main.go`

To POST a visit, use

```
curl -X 'POST' \
  'http://localhost:8080/visit' \
  -H 'accept: */*' \
  -H 'Content-Type: application/json' \
  -d '{
  "id": "d0e36d0e-3797-47f0-a740-433a78d23902",
  "url": "https://catropy.com"
}'
```

then fetch the counter with 

```
curl -X 'GET' 'http://localhost:8080/visit?u=https://catropy.com'
```

### Todos
* Use mutex
* write integration tests
* write unit tests
* discover code is set up wrong for unit testing, refactor code to make it more testable
* reconsider life choices
* tackle race condition
* add logging
