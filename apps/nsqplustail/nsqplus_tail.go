package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"internal/app"
	"internal/version"
	"nsqio/go-nsq"
	"strings"
)

var (
	flagSet     = flag.NewFlagSet("nsqplus-tail", flag.ExitOnError)
	showVersion = flagSet.Bool("version", false, "print version string")

	topic         = flagSet.String("topic", "", "NSQ topic")
	channel       = flagSet.String("channel", "", "NSQ channel")
	maxInFlight   = flagSet.Int("max-in-flight", 50, "max number of messages to allow in flight")
	totalMessages = flagSet.Int("n", 100, "total messages to show (will wait if starved)")

	nsqdTCPAddrs     = app.StringArray{}
	lookupdHTTPAddrs = app.StringArray{}
)

func init() {
	flagSet.Var(&nsqdTCPAddrs, "nsqd-tcp-address", "nsqd TCP address (may be given multiple times)")
	flagSet.Var(&lookupdHTTPAddrs, "lookupd-http-address", "lookupd HTTP address (may be given multiple times)")
}

type TailHandler struct {
	totalMessages int
	messagesShown int
}

func (th *TailHandler) HandleMessage(m *nsq.Message) error {
	th.messagesShown++
	_, err := os.Stdout.Write(m.Body)
	if err != nil {
		log.Fatalf("ERROR: failed to write to os.Stdout - %s", err)
	}
	_, err = os.Stdout.WriteString("\n")
	if err != nil {
		log.Fatalf("ERROR: failed to write to os.Stdout - %s", err)
	}
	if th.totalMessages > 0 && th.messagesShown >= th.totalMessages {
		os.Exit(0)
	}
	return nil
}

func main() {
	cfg := nsq.NewConfig()
	flagSet.Var(&nsq.ConfigFlag{cfg}, "consumer-opt", "option to passthrough to nsq.Consumer (may be given multiple times, http://godoc.org/nsqio/go-nsq#Config)")

	flagSet.Parse(os.Args[1:])

	if *showVersion {
		fmt.Printf("nsq_tail v%s\n", version.Binary)
		return
	}

	if *channel == "" {
		rand.Seed(time.Now().UnixNano())
		*channel = fmt.Sprintf("tail%06d#ephemeral", rand.Int()%999999)
	}

	if *topic == "" {
		log.Fatal("--topic is required")
	}

	if len(nsqdTCPAddrs) == 0 && len(lookupdHTTPAddrs) == 0 {
		log.Fatal("--nsqd-tcp-address or --lookupd-http-address required")
	}
	if len(nsqdTCPAddrs) > 0 && len(lookupdHTTPAddrs) > 0 {
		log.Fatal("use --nsqd-tcp-address or --lookupd-http-address not both")
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Don't ask for more messages than we want
	if *totalMessages > 0 && *totalMessages < *maxInFlight {
		*maxInFlight = *totalMessages
	}

	cfg.UserAgent = fmt.Sprintf("nsq_tail/%s go-nsq/%s", version.Binary, nsq.VERSION)
	cfg.MaxInFlight = *maxInFlight

	consumer, err := nsq.NewConsumer(*topic, *channel, cfg)
	if err != nil {
		log.Fatal(err)
	}

	consumer.AddHandler(&TailHandler{totalMessages: *totalMessages})

	err = consumer.ConnectToNSQDs(nsqdTCPAddrs)
	if err != nil {
		log.Fatal(err)
	}

	log.Output(2, fmt.Sprintf("lookupdHTTPAddrs length: %d, %v", len(lookupdHTTPAddrs), lookupdHTTPAddrs))
	err = consumer.ConnectToNSQLookupds(lookupdHTTPAddrs)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-consumer.StopChan:
			return
		case <-sigChan:
			consumer.Stop()
		}
	}
}
