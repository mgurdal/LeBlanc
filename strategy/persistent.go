package strategy

import (
	"fmt"
	"net"
	"sync"

	"github.com/google/uuid"
	"github.com/mgurdal/lb/service"
)

type persistent struct {
	Services []*service.Service
	mu       *sync.Mutex
	next     int
	Channels map[string]*service.Channel
}

func (p *persistent) ListServices() []*service.Service {
	return p.Services
}

// New returns persistent implementation(*persistent).
func NewPersistent(services []*service.Service) Strategy {

	return &persistent{
		Services: services,
		mu:       new(sync.Mutex),
		Channels: map[string]*service.Channel{},
	}
}

func (p *persistent) GetChannelByService(addr net.Addr) *service.Channel {
	for _, channel := range p.Channels {
		if channel.Dst.Addr.String() == addr.String() {
			return channel
		}
	}
	return nil
}

func (p *persistent) GetChannel(client *service.Client) *service.Channel {
	channel, ok := p.Channels[client.Addr.String()]

	if !ok {
		backend := p.Acquire(client)

		channel = &service.Channel{
			ID:  uuid.New(),
			Src: client,
			Dst: backend,
		}
		fmt.Printf("Registering %s\n", channel)
		p.Channels[client.Addr.String()] = channel
	} else {
		fmt.Printf("Using existing channel %s\n", channel)
	}

	return channel
}

func (p *persistent) Acquire(client *service.Client) *service.Service {
	s := p.Next()
	fmt.Printf("Selected %s for %s total: %d\n", s.Addr, client.Addr, len(p.Services))
	return s
}

// Next returns next address
func (p *persistent) Next() *service.Service {
	p.mu.Lock()
	services := p.ActiveServices()
	if len(services) == 0 {
		fmt.Println("No active service")
	}
	sc := services[p.next]
	p.next = (p.next + 1) % len(services)
	p.mu.Unlock()
	return sc
}

func (p *persistent) ActiveServices() []*service.Service {
	var ln int
	actives := make([]*service.Service, len(p.Services))
	for _, service := range p.Services {
		if service.Stat.Status == "HEALTHY" {
			actives[ln] = service
			ln++
		}

	}
	return actives[:ln]
}
