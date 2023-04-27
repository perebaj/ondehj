# ondehj




```go
go run cmd/ondehoje/main.go
```

getall
```bash
curl -XGET http://localhost:8000/event | jq 
```

create
```bash
curl -X POST -H "Content-Type: application/json" -d '{"Title":"Example Event","Description":"This is an example event.","Location":"New York City","StartTime":"2023-04-20T11:30:27.747223-04:00","EndTime":"2023-04-20T12:30:27.747223-04:00","InstagramPage":"example_event"}' -i http://localhost:8000/event
```

delete
```bash
curl -X DELETE -i http://localhost:8000/event/{id} 
```

update
```bash
curl -X PUT \
  -i http://localhost:8000/event/2 \
  -H 'Content-Type: application/json' \
  -d '{
        "title": "Novo título",
        "description": "Nova descrição",
        "location": "Nova localização",
        "start_time": "2023-05-01T10:00:00Z",
        "end_time": "2023-05-01T12:00:00Z",
        "instagram_page": "nova_pagina"
      }'

```

## Database approach

## Core Concepts 
* API
* Database
* Database Migration
* Structured Logs
* Dev Container environment
* Go Unit test
* Metrics