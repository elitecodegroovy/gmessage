package main

import (
	"runtime"
	"sync"
	"bytes"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"fmt"
	"net/url"
)

func declarePostFormData(request_url string, formData map[string][]string, wg *sync.WaitGroup){
	b := new(bytes.Buffer)
	//对form进行编码
	json.NewEncoder(b).Encode(formData)

	defer wg.Done()
	wg.Add(1)

	rsp, err := http.Post(request_url, "application/x-www-form-urlencoded", b)
	if err != nil {
		panic(err)
	}
	defer rsp.Body.Close()
	body_byte, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body_byte))
}

func postFormData(){
	maxProcs := runtime.NumCPU()
	runtime.GOMAXPROCS(maxProcs)

	var wg  sync.WaitGroup
	url := "https://test-01.biostime.us/pointstats/appeventspv"
	//map[string][]string
	js := `{"platform":1,"point_code":"1120100","campaign":"","url":"","created_time":1494317655854,"customer_id":"32737223","user_mark":"359320050689107","user_sourceid":"","mobile":"13500709760","mobile_info":"Nexus 6","app_version":"5.9.0","source_from":"","sku":"","spu":"","course_id":"","account_id":"","terminal_code":"","coupon_defid":"","expert_id":"","question_id":"","ip":"221.4.38.4","var_x":"","var_y":"","var_z":"","sign":"594fe4d7fa439e8e64dcce302a4cecdc"} `
	formData := url.Values{"pvdata": js,
	}

	for i := 0; i < 1; i++ {
		go declarePostFormData(url, formData , &wg)
	}
	wg.Wait()
}

func main() {
	//getOrderResp()
	postFormData()
}
