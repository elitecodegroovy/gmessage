package main

import (
	"flag"
	"log"
	"time"

	"github.com/elitecodegroovy/gmessage/gio"
)

func usage() {
	log.Fatalf("Usage: gio-req [-s server (%s)] <subject> <msg> \n", gio.DefaultURL)
}

func main() {
	var urls = flag.String("s", "nats://192.168.1.225:6222", "The gmessage server URLs (separated by comma)")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		usage()
	}

	nc, err := gio.Connect(*urls)
	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
	}
	defer nc.Close()
	subj, payload := args[0], []byte(args[1])

	msg, err := nc.Request(subj, []byte(payload), 100*time.Millisecond)
	if err != nil {
		if nc.LastError() != nil {
			log.Fatalf("Error in Request: %v\n", nc.LastError())
		}
		log.Fatalf("Error in Request: %v\n", err)
	}

	log.Printf("Published [%s] : '%s'\n", subj, payload)
	log.Printf("Received [%v] : '%s'\n", msg.Subject, string(msg.Data))
}
