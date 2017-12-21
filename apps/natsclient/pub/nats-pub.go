package main

import (
	"flag"
	"log"

	"nats-io/go-nats"
	"strconv"
	"time"
)

// NOTE: Use tls scheme for TLS, e.g. nats-pub -s tls://demo.nats.io:4443 foo hello
func usage() {
	log.Fatalf("Usage: nats-pub [-s server (%s)] <subject> <msg> \n", nats.DefaultURL)
}

func main() {
	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	//args := flag.Args()
	//if len(args) < 2 {
	//	usage()
	//}

	nc, err := nats.Connect(*urls)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	subj := "NATS"
	for i := 0; i < 10; i++ {
		currentTime := time.Now().Local()
		newFormat := currentTime.Format("2006-01-02 15:04:05.000")
		nc.Publish(subj, []byte("i++"+strconv.Itoa(i)+" time/"+newFormat))
		if i%100 == 0 {
			nc.Flush()
		}
	}
	nc.Flush()
	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Published [%s] :\n", subj)
	}
}
