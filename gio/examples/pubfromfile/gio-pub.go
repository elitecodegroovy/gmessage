package main

import (
	"flag"
	"log"

	"bufio"
	"fmt"
	"github.com/elitecodegroovy/gmessage/gio"
	"os"
	"strings"
)

func usage() {
	log.Fatalf("Usage: gio-pub [-s server (%s)] <subject> <msg> \n", gio.DefaultURL)
}

func doPublish(nc *gio.Conn, subj string) error {
	file, err := os.Open("file.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		if strings.TrimSpace(scanner.Text()) != "" {
			nc.Publish(subj, []byte(scanner.Text()))
		}
	}
	nc.Flush()
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return nil
}

func publishMsg(subj string) {
	var urls = flag.String("s", "nats://192.168.1.225:6222", "gmessage 服务器URL地址(使用逗号分隔多个地址)")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		usage()
	}

	nc, err := gio.Connect(*urls)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	doPublish(nc, subj)

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Published [%s] : \n", subj)
	}
}

func main() {
	//subj : topic message
	publishMsg("topic01")
}
