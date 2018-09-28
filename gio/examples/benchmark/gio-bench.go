package main

import (
	"flag"
	"log"
	"sync"
	"time"

	"fmt"
	"github.com/elitecodegroovy/gmessage/gio"
	"github.com/elitecodegroovy/gmessage/gio/bench"
	"io/ioutil"
	"strings"
)

// Some sane defaults
const (
	DefaultGMessage    = "gmessage://192.168.1.225:6222,gmessage://192.168.1.224:6222,gmessage://192.168.1.226:6222"
	DefaultNumMsgs     = 100
	DefaultNumPubs     = 1
	DefaultNumSubs     = 1
	DefaultMessageSize = 128
)

func usage() {
	log.Fatalf("Usage: gio-bench [-s server (%s)] [--tls] [-np NUM_PUBLISHERS] [-ns NUM_SUBSCRIBERS] [-n NUM_MSGS] [-ms MESSAGE_SIZE] [-csv csvfile] <subject>\n", DefaultGMessage)
}

var benchmark *bench.Benchmark

func startDoBenchmark() {
	var urls = flag.String("s", DefaultGMessage, "gmessage 服务器URL地址，多个地址以逗号分隔")
	var tls = flag.Bool("tls", false, "是否启用TLS安全连接")
	var numPubs = flag.Int("np", DefaultNumPubs, "并发发布者的个数")
	var numSubs = flag.Int("ns", DefaultNumSubs, "并发订阅者的个数")
	var numMsgs = flag.Int("n", DefaultNumMsgs, "发送的消息个数")
	var msgSize = flag.Int("ms", DefaultMessageSize, "消息的大小")
	var csvFile = flag.String("csv", "", "保存标准度量数据到csv文件")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		usage()
	}

	if *numMsgs <= 0 {
		log.Fatal("Number of messages should be greater than zero.")
	}

	// Setup the option block
	opts := gio.GetDefaultOptions()
	opts.Servers = strings.Split(*urls, ",")
	for i, s := range opts.Servers {
		opts.Servers[i] = strings.Trim(s, " ")
	}
	opts.Secure = *tls

	benchmark = bench.NewBenchmark("系统gmessage", *numSubs, *numPubs)

	var startwg sync.WaitGroup
	var donewg sync.WaitGroup

	donewg.Add(*numPubs + *numSubs)

	// Run Subscribers first
	startwg.Add(*numSubs)
	for i := 0; i < *numSubs; i++ {
		go runSubscriber(&startwg, &donewg, opts, *numMsgs, *msgSize)
	}
	startwg.Wait()

	// Now Publishers
	startwg.Add(*numPubs)
	pubCounts := bench.MsgsPerClient(*numMsgs, *numPubs)
	for i := 0; i < *numPubs; i++ {
		go runPublisher(&startwg, &donewg, opts, pubCounts[i], *msgSize)
	}

	log.Printf("开始启动基准度量测试 [消息量=%d, 消息大小=%d, 发布者=%d, 订阅者=%d]\n", *numMsgs, *msgSize, *numPubs, *numSubs)

	startwg.Wait()
	donewg.Wait()

	benchmark.Close()

	fmt.Print(benchmark.Report())

	if len(*csvFile) > 0 {
		csv := benchmark.CSV()
		ioutil.WriteFile(*csvFile, []byte(csv), 0644)
		fmt.Printf("保存标准度量数据到csv文件： %s\n", *csvFile)
	}
}

func runPublisher(startwg, donewg *sync.WaitGroup, opts gio.Options, numMsgs int, msgSize int) {
	nc, err := opts.Connect()
	if err != nil {
		log.Fatalf("无法连接: %v\n", err)
	}
	defer nc.Close()
	startwg.Done()

	args := flag.Args()
	subj := args[0]
	var msg []byte
	if msgSize > 0 {
		msg = make([]byte, msgSize)
	}

	start := time.Now()

	for i := 0; i < numMsgs; i++ {
		nc.Publish(subj, msg)
	}
	nc.Flush()
	benchmark.AddPubSample(bench.NewSample(numMsgs, msgSize, start, time.Now(), nc))

	donewg.Done()
}

func runSubscriber(startwg, donewg *sync.WaitGroup, opts gio.Options, numMsgs int, msgSize int) {
	nc, err := opts.Connect()
	if err != nil {
		log.Fatalf("无法连接: %v\n", err)
	}

	args := flag.Args()
	subj := args[0]

	received := 0
	start := time.Now()
	nc.Subscribe(subj, func(msg *gio.Msg) {
		received++
		if received >= numMsgs {
			benchmark.AddSubSample(bench.NewSample(numMsgs, msgSize, start, time.Now(), nc))
			donewg.Done()
			nc.Close()
		}
	})
	nc.Flush()
	startwg.Done()
}

func main() {
	startDoBenchmark()
}
