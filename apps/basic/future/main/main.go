package main

import (
	"net/http"
	"io/ioutil"
	"github.com/elitecodegroovy/gmessage/apps/basic/future"
	"fmt"
)

func callTask() {
	task := func()(r interface{}, err error){
		url := "http://ip.taobao.com/service/getIpInfo.php?ip=221.4.38.21"

		resp, err := http.Get(url)
		defer resp.Body.Close()
		if err != nil {
			return nil, err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return body, nil
	}

	f := promise.Start(task).OnSuccess(func(v interface{}) {
		fmt.Println("Success...")
		//...
	}).OnFailure(func(v interface{}) {
		fmt.Println("OnFailure...")
		//...
	}).OnComplete(func(v interface{}) {
		fmt.Println("OnComplete...")
		//...
	})
	if r, err := f.Get(); err == nil {
		fmt.Println("response : ", string(r.([]byte)))
	}else {
		fmt.Errorf("error %s", err.Error())
	}
}


func main() {
	//Test()
	callTask()
}
