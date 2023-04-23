# ondehj

```go
go run cmd/ondehoje/main.go
```


```bash
curl -X POST \ 
    -H "Content-Type: application/json" \
    -d '{"Title":"Example Event","Description":"This is an example event.","Location":"New York City","StartTime":"2023-04-20T11:30:27.747223-04:00" "EndTime":"2023-04-20T12:30:27.747223-04:00","InstagramPage":"example_event"}' \
    -i http://localhost:8080/event
```


```bash
curl -X DELETE -i http://localhost:8080/event/{id} 
```

## Database approach

## Core Concepts 
* API
* Database
* Database Migration
* Structured Logs
* Metrics
