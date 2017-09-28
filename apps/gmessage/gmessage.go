package main

import (
	"flag"
	"time"
	"log"
	"sync"
	"github.com/go-svc/svc"
)

var (
	flagSet = flag.NewFlagSet("gMessage", flag.ExitOnError)

	config      = flagSet.String("config", "", "path to config file")
	showVersion = flagSet.Bool("version", false, "print version string")
	verbose     = flagSet.Bool("verbose", false, "enable verbose logging")

	tcpAddress       = flagSet.String("tcp-address", "0.0.0.0:4160", "<addr>:<port> to listen on for TCP clients")
	httpAddress      = flagSet.String("http-address", "0.0.0.0:4161", "<addr>:<port> to listen on for HTTP clients")
	broadcastAddress = flagSet.String("broadcast-address", "", "address of this lookupd node, (default to the OS hostname)")

	inactiveProducerTimeout = flagSet.Duration("inactive-producer-timeout", 300*time.Second, "duration of time a producer will remain in the active list since its last ping")
	tombstoneLifetime       = flagSet.Duration("tombstone-lifetime", 45*time.Second, "duration of time a producer will remain tombstoned if registration remains")
)


// program implements svc.Service
type program struct {
	wg   sync.WaitGroup
	quit chan struct{}
}

func main() {
	prg := &program{}

	// Call svc.Run to start your program/service.
	if err := svc.Run(prg); err != nil {
		log.Fatal(err)
	}
}

func (p *program) Init(env svc.Environment) error {
	log.Printf("is win service? %v\n", env.IsWindowsService())
	return nil
}

func (p *program) Start() error {
	// The Start method must not block, or Windows may assume your service failed
	// to start. Launch a Goroutine here to do something interesting/blocking.

	p.quit = make(chan struct{})

	p.wg.Add(1)
	go func() {
		log.Println("Starting...")
		<-p.quit
		log.Println("Quit signal received...")
		p.wg.Done()
	}()

	return nil
}

func (p *program) Stop() error {
	// The Stop method is invoked by stopping the Windows service, or by pressing Ctrl+C on the console.
	// This method may block, but it's a good idea to finish quickly or your process may be killed by
	// Windows during a shutdown/reboot. As a general rule you shouldn't rely on graceful shutdown.

	log.Println("Stopping...")
	close(p.quit)
	p.wg.Wait()
	log.Println("Stopped.")
	return nil
}

