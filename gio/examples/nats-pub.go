
// +build ignore

package main

import (
	"flag"
	"log"

	"github.com/elitecodegroovy/gmessage/g-io"
)

// NOTE: Use tls scheme for TLS, e.g. nats-pub -s tls://demo.nats.io:4443 foo hello
func usage() {
	log.Fatalf("Usage: nats-pub [-s server (%s)] <subject> <msg> \n", nats.DefaultURL)
}

func main() {
	var urls = flag.String("s", "nats://192.168.1.225:6222", "The nats server URLs (separated by comma)")

	log.SetFlags(0)
	//flag.Usage = usage
	//flag.Parse()

	//args := flag.Args()
	//if len(args) < 2 {
	//	usage()
	//}

	nc, err := nats.Connect(*urls)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	subj, msg := "topic01", []byte("Msg Pub!")

	nc.Publish(subj, msg)
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Published [%s] : '%s'\n", subj, msg)
	}
}
