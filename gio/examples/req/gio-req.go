package main

import (
	"flag"
	"log"
	"time"

	"github.com/elitecodegroovy/gmessage/gio"
)

func usage() {
	log.Fatalf("用法: gio-req [-s server (%s)] <subject>（主题） <msg>（消息） \n", gio.DefaultURL)
}

func main() {
	var urls = flag.String("s",
		"gmessage://192.168.1.225:6222,gmessage://192.168.1.224:6222,gmessage://192.168.1.226:6222",
		"gmessage 服务器URL地址(使用逗号分隔多个地址)")

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
	defer nc.Close()
	subj, payload := "test01", []byte("10000")

	msg, err := nc.Request(subj, []byte(payload), 100*time.Millisecond)
	if err != nil {
		if nc.LastError() != nil {
			log.Fatalf("请求中的错误: %v\n", nc.LastError())
		}
		log.Fatalf("请求中的错误： %v\n", err)
	}

	log.Printf("完成 [%s] : '%s'\n", subj, payload)
	log.Printf("接受 [%v] : '%s'\n", msg.Subject, string(msg.Data))
}
