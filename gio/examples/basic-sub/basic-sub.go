package basic_sub

import (
	"flag"
	"fmt"
	"github.com/elitecodegroovy/gmessage/gio"
	"log"
)

func doSimpleSubsriber(subject string) {
	var urls = flag.String("s",
		"gmessage://192.168.1.225:6222,gmessage://192.168.1.224:6222,gmessage://192.168.1.226:6222",
		"gmessage 服务器URL地址(使用逗号分隔多个地址)")

	log.SetFlags(0)

	nc, err := gio.Connect(*urls)
	if err != nil {
		log.Fatalf("无法连接: %v\n", err)
	}

	defer nc.Close()
	// Simple Async Subscriber
	nc.Subscribe(subject, func(m *gio.Msg) {
		fmt.Printf("接受到的消息 %s\n", string(m.Data))
	})

}

func main() {
	doSimpleSubsriber("test01")
}
