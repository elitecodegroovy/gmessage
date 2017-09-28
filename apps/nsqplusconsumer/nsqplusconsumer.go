package main

import (
	"nsqplus"
	"encoding/json"
	"fmt"
	"time"
	c "nsqio/go-nsq"
	"strconv"
	"log"
	"os"
	"bytes"
	"net/http"
	"flag"
	"internal/app"
	"internal/version"
	"os/signal"
	"syscall"
	"toml"
	"go-options"
	"io/ioutil"
	"strings"
)

var (
	flagSet 		= flag.NewFlagSet("nsqplusconsumer", flag.ExitOnError)
	config      		= flagSet.String("config", "", "path to config file")

	showVersion 		= flagSet.Bool("version", false, "print version string")
	pushHttpAddress      	= flagSet.String("push-http-address", "http://10.50.115.14:8801/statistics/api/push", "<addr>:<port> to listen on for TCP clients")

	nsqlookupHTTPAddresses       = app.StringArray{}
)

func init() {
	flagSet.Var(&nsqlookupHTTPAddresses, "nsqlookup-http-addresses", "nsqlookup HTTP address (may be given multiple times)")
}



const(
	MAX_PLATFORM 	= 8
)

type MyHandler struct {
	q                *c.Consumer
	logger   	 *log.Logger
	messagesSent     int
	messagesReceived int
	messagesFailed   int
	PostUrl 	 string
}


func (h *MyHandler) LogFailedMessage(message *c.Message) {
	data := nsqplus.Event{}
	err := json.Unmarshal(message.Body, &data)
	if err != nil {
		h.logger.Output(2, fmt.Sprintf("error LogFailedMessage %#v", err))
	}
	h.logger.Output(2, fmt.Sprintf("[failed] json :%s ", message.Body))
	//h.q.Stop()
}

var totalCount int64 = 0
func (h *MyHandler) HandleMessage(message *c.Message) error {
	event := nsqplus.Event{}
	err := json.Unmarshal(message.Body, &event)
	if err != nil {
		return err
	}

	//h.logger.Output(2, fmt.Sprintf("req body: %#v ", event))
	h.messagesReceived++
	totalCount++
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(event)
	resp, err := http.Post(h.PostUrl, "application/json; charset=utf-8", b)
	if err != nil {
		h.logger.Output(2, fmt.Sprintf("response ERROR %s", err.Error()))
		return err
	}
	if resp.StatusCode == 200 { // OK
		_, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			fmt.Println("error ioutil.ReadAll"+ err2.Error())
		}

		//bodyString := string(bodyBytes)
		h.logger.Output(2, fmt.Sprintf("eq body: %#v ,total count : %d, isSuccess: true ", event, totalCount))
		//h.logger.Output(2, fmt.Sprintf("response body %s", bodyString))
		//h.logger.Output(2, fmt.Sprintf("response close %v", resp.Close))
	}else {
		h.logger.Output(3, fmt.Sprintf("eq body: %#v ,total count : %d, isSuccess: false ", event, totalCount))
		h.logger.Output(3, fmt.Sprintf("[error] request body %s", string(message.Body)))
	}


	return nil
}

func (h *MyHandler)reportStats(consumer *c.Consumer) error{
	stats := consumer.Stats()
	h.logger.Output(2, fmt.Sprintf(" received :%d, finished :%d, requeue: %d, connection NO : %d ",
		stats.MessagesReceived, stats.MessagesFinished, stats.MessagesRequeued, stats.Connections))
	return nil
}

func doConsuming(options *nsqplus.ConsumerOptions , cb func(c *c.Config)) {
	config := c.NewConfig()
	// so that the test can simulate reaching max requeues and a call to LogFailedMessage
	config.DefaultRequeueDelay = 0
	// so that the test wont timeout from backing off
	config.MaxBackoffDuration = time.Millisecond * 50
	if cb != nil {
		cb(config)
	}
	logger := log.New(os.Stderr, "[nsqplus-consumer] ", log.Ldate|log.Ltime|log.Lmicroseconds)
	logger.Output(2, fmt.Sprintf("taget http address: %s", options.PushHttpAddress))
	logger.Output(2, fmt.Sprintf("nsqpluslookup address: %s", strings.Join(options.NsqlookupHTTPAddress, ",")))
	for i:=0 ; i < MAX_PLATFORM; i++ {
		go listenTopic(options, config, i, logger)
	}
}

func listenTopic(options *nsqplus.ConsumerOptions, config  *c.Config, i int, logger *log.Logger) error{
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	logger.Output(1, fmt.Sprintf("%s %d", "start consumer ", i) )
	topicName := "biostime"+strconv.Itoa(i)
	consumer, err := c.NewConsumer(topicName, "hNh", config)
	if err != nil {
		logger.Fatal("stats NewConsumer  ", err)
		return err
	}
	consumer.SetLogger(logger, 1)

	h := &MyHandler{
		q: consumer,
		logger:		logger,
		PostUrl: 	options.PushHttpAddress,
	}
	consumer.AddHandler(h)

	h.messagesSent = 0

	err = consumer.ConnectToNSQLookupds(options.NsqlookupHTTPAddress)
	if err != nil {
		logger.Fatal("stats report 0 connections (should be > 0) ", err)
		return err
	}
	for {
		select {
		case <-consumer.StopChan:
		case <-sigChan:
			consumer.Stop()
		}
	}
	return nil
}

func postData(){
	event := nsqplus.Event{}
	js := `{"platform":1,"point_code":"1120100","campaign":"","url":"","created_time":1494317655854,"customer_id":"32737223","user_mark":"359320050689107","user_sourceid":"","mobile":"13500709760","mobile_info":"Nexus 6","app_version":"5.9.0","source_from":"","sku":"","spu":"","course_id":"","account_id":"","terminal_code":"","coupon_defid":"","expert_id":"","question_id":"","ip":"221.4.38.4","var_x":"","var_y":"","var_z":"","sign":"594fe4d7fa439e8e64dcce302a4cecdc"} `
	err := json.Unmarshal([]byte(js), &event)
	if err != nil {
		panic(err)
	}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(event)
	resp, _ := http.Post("http://localhost:8801/statistics/api/push", "application/json; charset=utf-8", b)
	defer resp.Body.Close()

	if resp.StatusCode == 200 { // OK
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			fmt.Println("error ioutil.ReadAll"+ err2.Error())
		}
		bodyString := string(bodyBytes)
		fmt.Println("response body : ", bodyString)
	}
}

func main(){
	flagSet.Parse(os.Args[1:])
	if *showVersion {
		fmt.Println(version.String("nsq-plus consumer"))
		return
	}

	exitChan := make(chan int)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		exitChan <- 1
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	//logic entrance
	var cfg map[string]interface{}
	if *config != "" {
		_, err := toml.DecodeFile(*config, &cfg)
		if err != nil {
			log.Fatalf("ERROR: failed to load config file %s - %s", *config, err)
		}
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	//NewConsumerOptions
	opts := nsqplus.NewConsumerOptions()
	options.Resolve(opts, flagSet, cfg)
	doConsuming(opts, func(c *c.Config) {
		c.Snappy = true
	})
	<-sigChan
}

