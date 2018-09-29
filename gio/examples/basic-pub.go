package main

import (
	"flag"
	"github.com/elitecodegroovy/gmessage/gio"
	"log"
	"strconv"
)

func doSimplePublish(subject, msg string) {
	var urls = flag.String("s",
		"gmessage://192.168.1.225:6222,gmessage://192.168.1.224:6222,gmessage://192.168.1.226:6222",
		"gmessage 服务器URL地址(使用逗号分隔多个地址)")

	log.SetFlags(0)

	nc, err := gio.Connect(*urls)
	if err != nil {
		log.Fatalf("无法连接: %v\n", err)
	}

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
	for i := 0; i < 100; i++ {
		doSimplePublish("test01", "msg:"+strconv.Itoa(i))
	}
}
