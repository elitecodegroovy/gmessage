package main

import (
	"flag"
	"log"
	"os"
	"runtime"

	"github.com/elitecodegroovy/gmessage/gio"
)

func usage() {
	log.Fatalf("用法: gio-qsub [-s server] [-t] <subject> <queue-group>\n")
}

func printMsg(m *gio.Msg, i int) {
	log.Printf("[#%d] Received on [%s] Queue[%s] Pid[%d]: '%s'\n", i, m.Subject, m.Sub.Queue, os.Getpid(), string(m.Data))
}

func main() {
	var urls = flag.String("s", "nats://192.168.1.225:6222",
		"gmessage 服务器URL地址(使用逗号分隔多个地址)")
	var showTime = flag.Bool("t", false, "显示时间戳")

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

	subj, queue, i := args[0], args[1], 0

	nc.QueueSubscribe(subj, queue, func(msg *gio.Msg) {
		i++
		printMsg(msg, i)
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
