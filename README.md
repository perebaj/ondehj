# ondehj

```go
go run api/main.go
```


```bash
curl -X POST -H "Content-Type: application/json" -d '{"Title":"Example Event","Description":"This is an example event.","Location":"New York City","StartTime":"2023-04-20T11:30:27.747223-04:00","EndTime":"2023-04-20T12:30:27.747223-04:00","InstagramPage":"example_event"}' http://localhost:3000/event
```