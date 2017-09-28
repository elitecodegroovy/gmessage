package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
	//	"bytes"
	//	"encoding/json"
	//	"net/url"
)

var httpGetResp = make(chan string)

func get(order string) {
	response, err := http.Get("http://weixin.mama100.cn/alipay/internationalizationwappay4website/query?out_trade_no=" + order + "&seq_no=0000000011111&mch_type=Swisse&from_system=MAMA100WEBSITE")
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	// response.Body() is a reader type. We have
	// to use ioutil.ReadAll() to read the data
	// in to a byte slice(string)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	httpGetResp <- fmt.Sprint("orderCode: %s : %s", order, string(body))
}

func declareGet(orderCode, transactionId string, wg *sync.WaitGroup) {
	defer wg.Done()
	//ZONGSHU
	uri := "http://192.168.234.204:42080/alipay/acquire/customs?seq_no=null&from_system=MAGELLAN&mchType=Swisse&out_request_no=" + orderCode + "&trade_no=" + transactionId + "&merchant_customs_code=1500002864&amount=33.00&customs_place=ZONGSHU&merchant_customs_name=广州市妈妈一百电子商务有限公司"
	response, err := http.Get(uri)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	// response.Body() is a reader type. We have
	// to use ioutil.ReadAll() to read the data
	// in to a byte slice(string)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("", uri, ", response result: " , string(body))
	httpGetResp <- fmt.Sprint("", uri, string(body))
}

func getOrderResp() {
	maxProcs := runtime.NumCPU()
	runtime.GOMAXPROCS(maxProcs)

	var wg sync.WaitGroup
	file, err := os.Open("./order.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		//get(scanner.Text())
		i++
		orderInfo := strings.Split(scanner.Text(), "\t")
		wg.Add(1)
		go declareGet(orderInfo[0], orderInfo[1], &wg)
		fmt.Println("i:", i)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	go func() {
		for m := range httpGetResp {
			fmt.Println(m)
		}
	}()
	wg.Wait()

}

func receiveMsg() {
	msg := make(chan string)
	go func() {
		time.Sleep(time.Second * 1)
		msg <- "goroutine 1"
	}()

	go func() {
		time.Sleep(time.Second * 2)
		msg <- "goroutine 2"
	}()
	go func() {
		time.Sleep(time.Second * 3)
		msg <- "goroutine 3"
	}()

	go func() {
		for i := range msg {
			fmt.Println("message :", i)
		}
	}()
	time.Sleep(time.Second * 4)
}

func receiveMsg2() {
	msg := make(chan string)
	done := make(chan bool)

	go func() {
		time.Sleep(time.Second * 1)
		msg <- "goroutine 3"
		done <- true
	}()
	go func() {
		time.Sleep(time.Second * 3)
		msg <- "goroutine 1"
		done <- true
	}()
	go func() {
		time.Sleep(time.Second * 2)
		msg <- "goroutine 2"
		done <- true
	}()

	go func() {
		time.Sleep(time.Millisecond * 1500)
		for i := range msg {
			fmt.Println("message :", i)
		}
	}()
	for j := 0; j < 3; j++ {
		<-done
	}
	time.Sleep(time.Second * 1)
}

func receiveMsg3() {
	msg := make(chan string)
	var wg sync.WaitGroup

	wg.Add(3)
	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 3)
		msg <- "goroutine 1"
	}()
	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 2)
		msg <- "goroutine 2"
	}()
	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 1)
		msg <- "goroutine 3"
	}()
	go func() {
		time.Sleep(time.Second * 3)
		for i := range msg {
			fmt.Println("message :", i)
		}
	}()
	wg.Wait()
}
