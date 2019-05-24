# API with CircuitBreaker
This is a simple api that has a single endpoint which calls out to `flakyservice` and returns the result. It uses the `sony/gobreaker` library to implement the circuit breaker pattern to reduce load on the backend service during prolonged disruptions.

### Usage
Start service
```
go run main.go
```

Access main endpoint
```
curl -v http://localhost:8081
```

### Results
Success
```
{"Message":"SUCCESS"}
```

Backend failed, but circuit breaker is still closed
```
{"Message":"Failed to connect to flakyservice: Backend failed"}
```

Circuit breaker is open and backend request was skipped
```
{"Message":"Failed to connect to flakyservice: circuit breaker is open"}
```