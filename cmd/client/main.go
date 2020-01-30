package main

import (
	"crypto/rand"
	"flag"
	"log"
	mrand "math/rand"
	"net"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

const (
	flushInterval = time.Duration(1) * time.Second
	UDPPacketSize = 1500
)

var address string
var bufferPool sync.Pool
var ops uint64 = 0
var total uint64 = 0
var flushTicker *time.Ticker
var nbWorkers int
var loading = true

func init() {
	flag.StringVar(&address, "addr", ":50007", "Address of the UDP server to test")
	flag.IntVar(&nbWorkers, "concurrency", runtime.NumCPU(), "Number of workers to run in parallel")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	bufferPool = sync.Pool{
		New: func() interface{} { return make([]byte, UDPPacketSize) },
	}
	load(nbWorkers)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			loading = false
			runtime.Gosched()
			atomic.AddUint64(&total, ops)
			log.Printf("Total ops %d", total)
			os.Exit(0)
		}
	}()

	flushTicker = time.NewTicker(flushInterval)
	for range flushTicker.C {
		log.Printf("Ops/s %f", float64(ops)/flushInterval.Seconds())
		atomic.AddUint64(&total, ops)
		atomic.StoreUint64(&ops, 0)
	}
}

func load(maxWorkers int) error {
	for i := 0; i < maxWorkers; i++ {
		go func() {
			for loading {
				mrand.Seed(time.Now().Unix())
				n := mrand.Intn(UDPPacketSize - 1)
				write(randBytes(n), n)
				time.Sleep(time.Duration(500) * time.Microsecond)
			}
		}()
	}
	return nil
}

func write(buf []byte, n int) {
	conn, err := net.Dial("udp", address)
	if err != nil {
		log.Printf("Error connecting to server: %s", err)
		return
	}
	defer conn.Close()
	defer func() { bufferPool.Put(buf) }()

	_, err = conn.Write(buf[0:n])
	if err != nil {
		log.Printf("Error sending to server: %s", err)
		return
	}
	atomic.AddUint64(&ops, 1)
}

func randBytes(n int) []byte {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes := bufferPool.Get().([]byte)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return bytes
}
