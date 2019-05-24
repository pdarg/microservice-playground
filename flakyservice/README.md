# Flaky Service
This is a simple HTTP service you can use to emulate a flaky backend service.

### Usage
Start service
```
go run main.go
```

Access main endpoint
```
curl -v http://localhost:8080
```

Enable failure mode. All requests with fail with a 500 status code
```
curl -v http://localhost:8080/stop
```

Enable success mode
```
curl -v http://localhost:8080/start
```