package httpclient

import (
	"net/http"
	"fmt"
	"io/ioutil"
)

func Get(url string)[]byte {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	// response.Body() is a reader type. We have
	// to use ioutil.ReadAll() to read the data
	// in to a byte slice(string)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Get url"+ url," with an error", err)
		return nil
	}
	return body
}

func GetWithParam(url string, params map[string]string)[]byte {
	var suffix string
	if params != nil && len(params) > 0 {
		for k, v := range params {
			suffix += k +"="+ v
		}
	}
	uri := fmt.Sprintf("%s?%s", url, suffix)
	return Get(uri)
}
