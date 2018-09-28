package main

import (
	"flag"
	"log"

	"github.com/elitecodegroovy/gmessage/gio"
)

func usage() {
	log.Fatalf("Usage: publish [-s server (%s) -t topic_name -c topic's_content]  \n", gio.DefaultURL)
}

func publishMsg() {
	var urls = flag.String("s", "gmessage://192.168.1.225:6222", "The gmessage server URLs (separated by comma)")
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

	subj, reply := args[0], args[1]

	nc.Publish(subj, []byte(reply))
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Published [%s] : %s\n", subj, reply)
	}
}

func main() {
	publishMsg()
}
