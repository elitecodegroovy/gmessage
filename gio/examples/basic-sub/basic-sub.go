package main

import (
	"flag"
	"github.com/elitecodegroovy/gmessage/gio"
	"log"
	"runtime"
	"os"
)


func printMsg(m *gio.Msg, i int) {
	log.Printf("[#%d] Received on [%s] Queue[%s] Pid[%d]: '%s'\n", i, m.Subject, m.Sub.Queue, os.Getpid(), string(m.Data))
}

func doSimpleSubsriber() {
	var urls = flag.String("s",
		"gmessage://192.168.1.225:6222",
		"gmessage 服务器URL地址(使用逗号分隔多个地址)")

	log.SetFlags(0)
	//flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		log.Fatalf("请输入主题参数： baisc-sub subjectName \n")
	}
	subject := args[0]
	nc, err := gio.Connect(*urls)
	if err != nil {
		log.Fatalf("无法连接: %v\n", err)
	}

	i := 0
	// Simple Async Subscriber
	nc.Subscribe(subject, func(msg *gio.Msg) {
		i++
		printMsg(msg, i)
	})
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.SetFlags(log.LstdFlags)
	log.Printf("Listening on [%s]\n", subject)



	runtime.Goexit()
}

func main() {
	doSimpleSubsriber()
}
