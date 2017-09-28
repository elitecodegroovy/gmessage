package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type GeoIP struct {
	// The right side is the name of the JSON variable
	Ip          string  `json:"ip"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	RegionCode  string  `json:"region_code"`
	RegionName  string  `json:"region_name"`
	City        string  `json:"city"`
	Zipcode     string  `json:"zipcode"`
	Lat         float32 `json:"latitude"`
	Lon         float32 `json:"longitude"`
	MetroCode   int     `json:"metro_code"`
	AreaCode    int     `json:"area_code"`
}

type IpInfo struct {
	Ret 		int  	`json:"ret"`
	Atart           int  	`json:"start"`
	End 		int  	`json:"end"`
	Country  	string 	`json:"country"`
	Province 	string  `json:"province"`
	City 		string 	`json:"city"`
	District 	string  `json:"district"`
	Isp 		string  `json:"isp"`
	T_type           string `json:"type"`
	Desc 		string  `json:"desc"`
}

var (
	address  string
	err      error
	geo      GeoIP
	ipInfo   IpInfo
	response *http.Response
	body     []byte
)

func showIpInfo(ip string) IpInfo {
	go func(){
		time.Sleep(time.Second * 1)
	}()

	response, err = http.Get("http://int.dpool.sina.com.cn/iplookup/iplookup.php?format=json&ip=" + ip)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	// response.Body() is a reader type. We have
	// to use ioutil.ReadAll() to read the data
	// in to a byte slice(string)
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("resp body \n", string(body))
	// Unmarshal the JSON byte slice to a GeoIP struct

	err = json.Unmarshal(body, &ipInfo)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("ipInfo country :", ipInfo.Country )
	fmt.Println("ipInfo province", ipInfo.Province)
	fmt.Println("ipInfo city", ipInfo.City)
	fmt.Println("ip info :", ipInfo)
	return ipInfo
}

func goLocation(){
	// Provide a domain name or IP address
	address = "www.ifeng.com"
	// address = "2600:3c00::f03c:91ff:fe98:c0f5"
	// address = "45.79.8.237"

	// Use freegeoip.net to get a JSON response
	// There is also /xml/ and /csv/ formats available
	response, err = http.Get("https://freegeoip.net/json/" + address)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	// response.Body() is a reader type. We have
	// to use ioutil.ReadAll() to read the data
	// in to a byte slice(string)
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("resp body \n", string(body))

	// Unmarshal the JSON byte slice to a GeoIP struct
	err = json.Unmarshal(body, &geo)
	if err != nil {
		fmt.Println(err)
	}

	// Everything accessible in struct now
	fmt.Println("\n==== IP Geolocation Info ====\n")
	fmt.Println("IP address:\t", geo.Ip)
	fmt.Println("Country Code:\t", geo.CountryCode)
	fmt.Println("Country Name:\t", geo.CountryName)
	fmt.Println("Zip Code:\t", geo.Zipcode)
	fmt.Println("Latitude:\t", geo.Lat)
	fmt.Println("Longitude:\t", geo.Lon)
	fmt.Println("Metro Code:\t", geo.MetroCode)
	fmt.Println("Area Code:\t", geo.AreaCode)
	showIpInfo(geo.Ip)
}
//func main() {
//	goLocation()
//}