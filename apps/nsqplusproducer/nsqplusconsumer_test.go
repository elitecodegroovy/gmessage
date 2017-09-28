package main

import (
	"testing"
	"nsqplus"
	"encoding/json"
	"fmt"
	"time"
	c "nsqio/go-nsq"
)

type MyHandler struct {
	t                *testing.T
	q                *c.Consumer
	messagesSent     int
	messagesReceived int
	messagesFailed   int
}


func (h *MyHandler) LogFailedMessage(message *c.Message) {
	h.messagesFailed++
	h.q.Stop()
}

var totalCount int64 = 0
func (h *MyHandler) HandleMessage(message *c.Message) error {
	data := nsqplus.Event{}
	err := json.Unmarshal(message.Body, &data)
	if err != nil {
		return err
	}

	fmt.Sprintf("req body: %#v ", data)
	h.messagesReceived++
	totalCount++
	fmt.Println("total count ", totalCount)
	return nil
}


func TestMyConsumer(t *testing.T){
	doMyConsumerTest(t, func(c *c.Config) {
		c.Snappy = true
	})
}

func doMyConsumerTest(t *testing.T, cb func(c *c.Config)) {
	config := c.NewConfig()
	// so that the test can simulate reaching max requeues and a call to LogFailedMessage
	config.DefaultRequeueDelay = 0
	// so that the test wont timeout from backing off
	config.MaxBackoffDuration = time.Millisecond * 50
	if cb != nil {
		cb(config)
	}
	topicName := "biostime-test"
	q, _ := c.NewConsumer(topicName, "ch", config)
	// q.SetLogger(nullLogger, LogLevelInfo)

	h := &MyHandler{
		t: t,
		q: q,
	}
	q.AddHandler(h)

	h.messagesSent = 4

	addr := "10.50.115.16:4150"
	err := q.ConnectToNSQLookupd("10.50.115.16:4161")
	if err != nil {
		t.Fatal(err)
	}

	stats := q.Stats()
	if stats.Connections == 0 {
		t.Fatal("stats report 0 connections (should be > 0)")
	}

	err = q.ConnectToNSQD(addr)
	if err == nil {
		t.Fatal("should not be able to connect to the same NSQ twice")
	}

	<-q.StopChan

	stats = q.Stats()
	if stats.Connections != 0 {
		t.Fatalf("stats report %d active connections (should be 0)", stats.Connections)
	}

	stats = q.Stats()
	fmt.Sprintf("message received %d, finished %d", stats.MessagesReceived, stats.MessagesFinished)

	fmt.Sprintf("handler messagesReceived %d, sent %d by os ,failed %d", h.messagesReceived, h.messagesSent, h.messagesFailed)
}
