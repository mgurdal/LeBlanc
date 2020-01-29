package strategy

import "github.com/mgurdal/lb/service"

type Strategy interface {
	Acquire(client *service.Client) *service.Service
	GetChannel(client *service.Client) *service.Channel
	ListServices() []*service.Service
	ActiveServices() []*service.Service
}
