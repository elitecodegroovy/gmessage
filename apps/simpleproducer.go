package main

import (
	"io/ioutil"
	"log"
	"nsq" //default ,you should import the package full path 'github.com/bitly/go-nsq' with the command `go get github.com/bitly/go-nsq`
	//"strconv"
	//	"strconv"
)

var nullLogger = log.New(ioutil.Discard, "", log.LstdFlags)

func sendMsg(message string) {
	//init default config
	config := nsq.NewConfig()

	w, _ := nsq.NewProducer("192.168.234.3:4150,192.168.234.27:4150,192.168.234.29:4150", config)
	err := w.Ping()
	if err != nil {
		//192.168.2.117:4150,192.168.2.68:4150
		log.Fatalln("error ping 10.50.115.16:4150", err)
		// switch the second nsq. You can use nginx or HAProxy for HA.
		w, _ = nsq.NewProducer("192.168.2.68:4150", config)
	}
	w.SetLogger(nullLogger, nsq.LogLevelInfo)

	err2 := w.Publish("a-test", []byte(message))
	if err2 != nil {
		log.Panic("Could not connect nsq")
	}

	w.Stop()
}

//func main() {
//	for i := 0; i < 2; i ++ {
//		sendMsg("msg index "+ strconv.Itoa(i + 10000))
//	}
//}
