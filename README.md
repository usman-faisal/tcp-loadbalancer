# TCP Load Balancer

A Go TCP load balancer with pluggable scheduling algorithms, connection tracking, and basic health checks.

![alt text](/assets/images/diagram.png)

## How it works

Backends are registered in a scheduler chosen by the `algorithm` config value. The load balancer accepts a connection, picks a backend, dials it, and proxies traffic bidirectionally. When a dial fails, the backend is marked unhealthy and skipped until it recovers. Connections are tracked and released when either side closes.

### Scheduling algorithms

- **round-robin** — picks backends in a fixed rotation, spreading connections evenly regardless of load.
- **least-conn** — keeps backends in a min-heap ordered by active connection count and routes each new connection to the least busy backend.

## Config

Backends and algorithm are set in `config.yaml`:

```yaml
port: 80
backend_list:
  - localhost:8000
  - localhost:8001
algorithm: round-robin
```

`algorithm` accepts either `round-robin` or `least-conn`.

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