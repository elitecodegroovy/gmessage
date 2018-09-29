package main

import (
	"flag"
	"github.com/elitecodegroovy/gmessage/gio"
	"log"
	"runtime"
)

func usage() {
	log.Fatalf("用法: gio-rply [-s server][-t] <subject> <response>\n")
}

func printMsg(m *gio.Msg, i int) {
	log.Printf("[#%d] 接受 [%s]: '%s'\n", i, m.Subject, string(m.Data))
}

func doReply() {
	var urls = flag.String("s",
		"gmessage://192.168.1.225:6222,gmessage://192.168.1.224:6222,gmessage://192.168.1.226:6222",
		"gmessage 服务器URL地址，多个地址以逗号分隔")
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
		log.Fatalf("无法连接: %v\n", err)
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

	log.Printf("监听主题 [%s]\n", subj)
	if *showTime {
		log.SetFlags(log.LstdFlags)
	}

	runtime.Goexit()
}

func main() {
	doReply()
}
