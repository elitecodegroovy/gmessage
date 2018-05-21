package httpclient

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"log"
	"time"
)

var myTransport http.RoundTripper = &http.Transport{
	Proxy:                 http.ProxyFromEnvironment,
	ResponseHeaderTimeout: time.Second * 2,
}

func DoHTTPReq(uri string) ([]byte, error) {
	fmt.Printf("HTTPDownload From: %s.\n", uri)
	var client = &http.Client{Transport: myTransport}
	res, err := client.Get(uri)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ReadFile: Size of download: %d\n", len(d))
	return d, err
}

func WriteFile(dst string, d []byte) error {
	fmt.Printf("WriteFile: Size of download: %d\n", len(d))
	err := ioutil.WriteFile(dst, d, 0444)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func DownloadImg(uri string, dst string) {
	fmt.Printf("下载文件的URL地址: %s.\n", uri)
	if d, err := DoHTTPReq(uri); err == nil {
		fmt.Printf("downloaded %s.\n", uri)
		if WriteFile(dst, d) == nil {
			fmt.Printf("saved %s as %s\n", uri, dst)
		}
	}
}


