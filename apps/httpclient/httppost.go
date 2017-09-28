package httpclient

import (
	"net/http"
	"bytes"
	"fmt"
	"io/ioutil"
)

func PostWithJsonBody(url string, js []byte)([]byte, error){
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(js))
	req.Header.Set("X-Custom-Header", "POST")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}
	return body , nil
}
