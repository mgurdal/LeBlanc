package service

import (
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	ID   uuid.UUID
	Addr net.Addr
	Conn net.PacketConn
	Stat *Stat
}

type Client struct {
	ID   uuid.UUID
	Addr net.Addr
	Conn net.PacketConn
}

type Stat struct {
	Status  string
	Latency int
}

func HealtCheck(service *Service) Stat {
	return Stat{}
}

type Channel struct {
	ID  uuid.UUID
	Src *Client
	Dst *Service
}

func (ch *Channel) String() string {
	return fmt.Sprintf("Channel(%s -> %s)", ch.Src.Addr, ch.Dst.Addr)
}

// Push sends the backend data to client
func (ch *Channel) Push(readBuffer []byte) {
	timeout := time.Second * 10
	time.Sleep(time.Second * 2)
	// server write
	deadline := time.Now().Add(timeout)
	err := ch.Src.Conn.SetWriteDeadline(deadline)
	if err != nil {
		return
	}

	n, err := ch.Src.Conn.WriteTo(readBuffer, ch.Src.Addr)
	if err != nil {
		return
	}
	fmt.Printf("packet-written: bytes=%d to=%s\n", n, ch.Src.Addr)
	return
}
