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
	var urls = flag.String("s", "nats://192.168.1.225:6222",
		"gmessage 服务器URL地址(使用逗号分隔多个地址)")
	var topic = flag.String("t", "test01",
		"主题名称定义，默认‘test01’")
	var content = flag.String("c", "test...",
		"主题的内容，默认'test...' ")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) > 4 {
		usage()
	}

	nc, err := gio.Connect(*urls)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	nc.Publish(*topic, []byte(*content))
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Published [%s] : %s\n", *topic, *content)
	}
}

func main() {
	publishMsg()
}
