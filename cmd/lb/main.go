package main

import (
	"net"

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
				Port: 500121,
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
	}

	strategy := strategy.NewPersistent(services)
	lb := lb.LB{Strategy: strategy}

	addr := ":50007"
	lb.Listen(addr)
}
