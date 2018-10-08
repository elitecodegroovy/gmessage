package main

import (
	"log"
	"github.com/elitecodegroovy/gmessage/gio"
	"flag"
	"runtime"
)

// NOTE: Use tls scheme for TLS, e.g. gio-rply -s tls://demo.nats.io:4443 foo hello
func usage() {
	log.Fatalf("Usage: gio-rply [-s server][-t] <subject> <response>\n")
}

func printMsg(m *gio.Msg, i int) {
	log.Printf("[#%d] Received on [%s]: '%s'\n", i, m.Subject, string(m.Data))
}

func doReply(){
	var urls = flag.String("s", "nats://192.168.1.225:6222", "The gio server URLs (separated by comma)")
	var showTime = flag.Bool("t", false, "Display timestamps")

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

	subj, reply, i := args[0], args[1], 0

	nc.Subscribe(subj, func(msg *gio.Msg) {
		i++
		printMsg(msg, i)
		nc.Publish(msg.Reply, []byte(reply))
	})
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on [%s]\n", subj)
	if *showTime {
		log.SetFlags(log.LstdFlags)
	}

	runtime.Goexit()
}

func main() {
	doReply()
}
