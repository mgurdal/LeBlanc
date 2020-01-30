# LeBlanc
L4 Load Balancer


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
