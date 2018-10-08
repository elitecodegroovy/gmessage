package main

import (
	"flag"
	"github.com/elitecodegroovy/gmessage/gio"
	"log"
)

func doSimplePublish() {
	var urls = flag.String("s",
		"gmessage://192.168.1.225:6222,gmessage://192.168.1.224:6222,gmessage://192.168.1.226:6222",
		"gmessage 服务器URL地址(使用逗号分隔多个地址)")
	log.SetFlags(0)
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		log.Fatalf("请输入主题参数： baisc-sub subjectName \n")
	}

	nc, err := gio.Connect(*urls)
	if err != nil {
		log.Fatalf("无法连接: %v\n", err)
	}

	subject, msg := args[0], args[1]

	// 简单发布消息
	nc.Publish(subject, []byte(msg))

	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Published [%s] : %s\n", subject, msg)
	}
}

func main() {
	doSimplePublish()
}
