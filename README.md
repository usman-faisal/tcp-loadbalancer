# TCP Load Balancer

A minimal Go TCP load balancer that routes connections to the backend with the fewest active connections.

![alt text](/assets/images/diagram.png)

## How it works

Backends are tracked in a min-heap ordered by active connection count. Each incoming connection goes to the least busy backend, which gets dialed and proxied bidirectionally to the client. When the connection closes, the count drops and the heap updates so the next request goes wherever's least loaded.

## Config

Backends are set in `config.yaml`:

```yaml
port: 80
backend_list:
  - localhost:8000
  - localhost:8001
```

## Running it

Start the sample backends:

```bash
go run cmd/server/main.go -port 8000
go run cmd/server/main.go -port 8001
```

Start the load balancer:

```bash
go run ./cmd/loadbalancer
```

Test it:

```bash
curl http://localhost
```