package strategy

import (
	"net"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"github.com/mgurdal/lb/service"
)

type roundrobin struct {
	Channels []*service.Channel
	Services []*service.Service
	mu       *sync.Mutex
	next     int
}

// New returns RoundRobin implementation(*roundrobin).
func NewRobin(services []*service.Service) Strategy {
	return &roundrobin{
		Services: services,
		mu:       new(sync.Mutex),
	}
}

func (r *roundrobin) GetChannelByService(addr net.Addr) *service.Channel {
	for _, channel := range r.Channels {
		if channel.Dst.Addr.String() == addr.String() {
			return channel
		}
	}
	return nil
}

func (r *roundrobin) GetChannel(client *service.Client) *service.Channel {
	backend := r.Acquire(client)
	channel := &service.Channel{
		ID:  uuid.New(),
		Src: client,
		Dst: backend,
	}
	return channel
}

// Next returns next address
func (r *roundrobin) Next() *service.Service {
	r.mu.Lock()
	services := r.ActiveServices()
	sc := services[r.next]
	r.next = (r.next + 1) % len(services)
	r.mu.Unlock()
	return sc
}

func (rr *roundrobin) Acquire(client *service.Client) *service.Service {

	s := rr.Next()
	// fmt.Printf("Selected %s for %s total: %d", s.Addr, client.Addr, len(rr.Services))
	return s
}

func (r *roundrobin) ListServices() []*service.Service {
	return r.Services
}

func (r *roundrobin) ActiveServices() []*service.Service {
	var ln int
	actives := make([]*service.Service, len(r.Services))
	for _, service := range r.Services {
		if service.Stat.Status == "HEALTHY" {
			actives[ln] = service
			ln++
		}

	}
	return actives[:ln]
}

func (r *roundrobin) Stats() map[string]string {
	return map[string]string{
		"Channels": strconv.Itoa(len(r.Channels)),
		"Services": strconv.Itoa(len(r.ActiveServices())),
	}
}
