// +build !race

package server

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nats-io/go-nats"
)

// IMPORTANT: Tests in this file are not executed when running with the -race flag.

func TestAvoidSlowConsumerBigMessages(t *testing.T) {
	opts := DefaultOptions() // Use defaults to make sure they avoid pending slow consumer.
	s := RunServer(opts)
	defer s.Shutdown()

	nc1, err := gio.Connect(fmt.Sprintf("nats://%s:%d", opts.Host, opts.Port))
	if err != nil {
		t.Fatalf("Error on connect: %v", err)
	}
	defer nc1.Close()

	nc2, err := gio.Connect(fmt.Sprintf("nats://%s:%d", opts.Host, opts.Port))
	if err != nil {
		t.Fatalf("Error on connect: %v", err)
	}
	defer nc2.Close()

	data := make([]byte, 1024*1024) // 1MB payload
	rand.Read(data)

	expected := int32(500)
	received := int32(0)

	done := make(chan bool)

	// Create Subscription.
	nc1.Subscribe("slow.consumer", func(m *gio.Msg) {
		// Just eat it so that we are not measuring
		// code time, just delivery.
		atomic.AddInt32(&received, 1)
		if received >= expected {
			done <- true
		}
	})

	// Create Error handler
	nc1.SetErrorHandler(func(c *gio.Conn, s *gio.Subscription, err error) {
		t.Fatalf("Received an error on the subscription's connection: %v\n", err)
	})

	nc1.Flush()

	for i := 0; i < int(expected); i++ {
		nc2.Publish("slow.consumer", data)
	}
	nc2.Flush()

	select {
	case <-done:
		return
	case <-time.After(10 * time.Second):
		r := atomic.LoadInt32(&received)
		if s.NumSlowConsumers() > 0 {
			t.Fatalf("Did not receive all large messages due to slow consumer status: %d of %d", r, expected)
		}
		t.Fatalf("Failed to receive all large messages: %d of %d\n", r, expected)
	}
}
