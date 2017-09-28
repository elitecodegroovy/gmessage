package main

import (
	"log"
	"nsq"
	"fmt"
)


func doSimpleConsumerTask(){

	config := nsq.NewConfig()
	q, _ := nsq.NewConsumer("a-test", "ch", config)
	q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		log.Printf("message: %v", string(message.Body))
		message.Finish()
		return nil
	}))
	lookupAddr := []string {
		"192.168.234.77:4161",
		"192.168.234.36:4161",
		"192.168.234.39:4161",
	}
	err := q.ConnectToNSQLookupds(lookupAddr)
	if err != nil {
		log.Panic("Could not connect")
	}
	<-q.StopChan

	stats := q.Stats()
	fmt.Sprintf("message received %d, finished %d", stats.MessagesReceived, stats.MessagesFinished)
}
//func main(){
//	doSimpleConsumerTask()
//}
