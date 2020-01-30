package lb

import (
	"fmt"
	"net"
	"time"

	"github.com/mgurdal/lb/service"
	"github.com/mgurdal/lb/strategy"
)

var (
	Info = Teal
	Warn = Yellow
	Fata = Red
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

var (
	Black   = Color("\033[1;30m%s\033[0m")
	Red     = Color("\033[1;31m%s\033[0m")
	Green   = Color("\033[1;32m%s\033[0m")
	Yellow  = Color("\033[1;33m%s\033[0m")
	Purple  = Color("\033[1;34m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
	White   = Color("\033[1;37m%s\033[0m")
)

const (
	HEALTHY             = "HEALTHY"
	DEAD                = "DEAD"
	MAX_UDP_PACKET_SIZE = 1500

	BANNER = `
	▄█          ▄████████ ▀█████████▄   ▄█          ▄████████ ███▄▄▄▄    ▄████████ 
	███         ███    ███   ███    ███ ███         ███    ███ ███▀▀▀██▄ ███    ███ 
	███         ███    █▀    ███    ███ ███         ███    ███ ███   ███ ███    █▀  
	███        ▄███▄▄▄      ▄███▄▄▄██▀  ███         ███    ███ ███   ███ ███        
	███       ▀▀███▀▀▀     ▀▀███▀▀▀██▄  ███       ▀███████████ ███   ███ ███        
	███         ███    █▄    ███    ██▄ ███         ███    ███ ███   ███ ███    █▄  
	███▌    ▄   ███    ███   ███    ███ ███▌    ▄   ███    ███ ███   ███ ███    ███ 
	█████▄▄██   ██████████ ▄█████████▀  █████▄▄██   ███    █▀   ▀█   █▀  ████████▀  
	▀                                   ▀                                           
	`
)

type LB struct {
	Strategy strategy.Strategy
	Active   int
	Total    int
	Latency  time.Duration
}

// Route selects a server based on the
// pre-defined load balancing strategy
// and forwards the dataflow to the server
func (lb *LB) Route(conn net.PacketConn) {
	defer conn.Close()
	for {

		readBuffer := make([]byte, MAX_UDP_PACKET_SIZE)

		// Read client input
		n, clientAddr, err := conn.ReadFrom(readBuffer)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Recieved a message from one of the backends
		channel := lb.Strategy.GetChannelByService(clientAddr)
		if channel != nil {
			fmt.Printf("Received from server %s (%s)\n", channel.Dst.Addr, clientAddr)
			go channel.ReSend(readBuffer)
		} else {
			fmt.Printf("packet-received: bytes=%d from=%s\n", n, clientAddr.String())

			client := &service.Client{
				Addr: clientAddr,
				Conn: conn,
			}

			channel := lb.Strategy.GetChannel(client)
			// TODO: check service availability
			go channel.Push(readBuffer[:n])
		}

	}

}

// Find returns the proper service for the given
// client
func (lb *LB) Find(client *service.Client) *service.Service {
	// Fresh connection

	return lb.Strategy.Acquire(client)

}

// Discover dials UDP connection to all registered servers.
// Records server conditions under stats.
func (lb *LB) Discover() {
	for _, service := range lb.Strategy.ListServices() {

		serverConn, err := net.DialUDP("udp", nil, service.Addr.(*net.UDPAddr))
		if err != nil {
			service.Stat.Status = DEAD
		} else if serverConn.RemoteAddr() == nil {
			service.Stat.Status = DEAD
		} else {
			service.Stat.Status = HEALTHY
			service.Conn = serverConn
		}
		fmt.Println(service.Addr, service.Stat.Status)
	}
}

// Listen starts the load balancer at given address
func (lb *LB) Listen(addr string) {
	lbConn, err := net.ListenPacket("udp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(Black(BANNER))
	fmt.Println(Info(fmt.Sprintf("Serving at: %s", lbConn.LocalAddr())))
	lb.Discover()

	lb.Route(lbConn)

}
