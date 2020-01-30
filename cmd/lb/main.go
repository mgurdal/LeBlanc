package main

import (
	"net"
	"sync"

	"github.com/google/uuid"
	"github.com/mgurdal/lb/lb"
	"github.com/mgurdal/lb/service"
	"github.com/mgurdal/lb/strategy"
)

func main() {
	services := []*service.Service{
		&service.Service{
			ID: uuid.New(),
			Addr: &net.UDPAddr{
				IP:   net.IPv4(127, 0, 0, 1),
				Port: 50004,
			},
			Stat: &service.Stat{},
		},
		&service.Service{
			ID: uuid.New(),
			Addr: &net.UDPAddr{
				IP:   net.IPv4(127, 0, 0, 1),
				Port: 50005,
			},
			Stat: &service.Stat{},
		},
		&service.Service{
			ID: uuid.New(),
			Addr: &net.UDPAddr{
				IP:   net.IPv4(127, 0, 0, 1),
				Port: 50009,
			},
			Stat: &service.Stat{},
		},
	}

	strategy := strategy.NewRobin(services)
	lb := lb.LB{Strategy: strategy, Mu: new(sync.Mutex)}

	addr := ":50007"
	lb.Listen(addr)
}
