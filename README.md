# LeBlanc
TCP Load Balancer


## ZERO CODE
```yaml
version: '3'

services:
  lb:
    image: leblanc
    restart: unless-stopped
    ports:
      - "60001:60001/udp"
    environments:
      - LB_CONFIG_FILE=lb.config
  server:
    image: leblanc/slow-server
    ports:
        - "60002-60010:60001"
```
```bash
docker-compose up
```


# Manual Setup
```go

func main() {

    services := []string{
        "udp://10.12.32.131:50001",
        "udp://10.12.32.132:50001",
    }

    strategy := strategy.NewPersistent(services)
    lb := lb.LB{Strategy: strategy}

    lb.Listen(":50007")
}
```

## Monitor

```go

lb.Services()

| Server                 | status    | I/O/Sec   |   N   |
|------------------------|-----------|-----------|-------|
| udp 10.12.32.131 50001 | healthy   | 421/109   |  102  |
| udp 10.12.32.132 50001 | dead      |   0/0     |  0    |
| udp 10.12.32.133 50001 | busy      | 102/100   |  30   |
|                        |           |           |       |


lb.Clients()

| Client                 | status    |  Server            | 
|------------------------|-----------|--------------------|
| udp 127.0.0.1    54135 | connected | 10.12.32.131 50001 |
|                        |           |                    |

```
