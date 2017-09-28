package main

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"
	"strconv"
	"log"
	"strings"
	"sync"
	"flag"
	"os"
	"unsafe"
	"net"
)

var a uint64 = 0

func getName(params ...interface{})string {
	var paramSlice []string
	for _, param := range params {
		switch v := param.(type){
		case string :
			paramSlice = append(paramSlice, v)
		case int :
			paramSlice = append(paramSlice , strconv.Itoa(v))
		default:
			log.Fatalln(" params has error type", v)
		}
	}
	return strings.Join(paramSlice, " ")
}

func test1(){
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println(runtime.NumCPU(), runtime.GOMAXPROCS(0))

	go func() {
		for {
			atomic.AddUint64(&a, uint64(1))
		}
	}()

	for {
		val := atomic.LoadUint64(&a)
		fmt.Println(val)
		time.Sleep(time.Second)
	}
}

//Best way to get single pattern object
type singleton struct {
}

var instance *singleton
var once sync.Once

func GetInstance() *singleton {
	once.Do(func() {
		instance = &singleton{}
	})
	return instance
}

func fibonacci(n int , c chan int){
	x , y := 1, 1
	for i :=0; i < n; i++ {
		c <- x
		x, y = y, x + y
	}
	close(c)
	isClosed, ok := <- c
	if ok {
		fmt.Println("closed", isClosed)
	}
}

func testFibonacci(){
	c := make(chan int , 10)

	go fibonacci(cap(c), c)
	for j := range c {
		fmt.Println(" ", j)
	}
	testSelectWaitNSleep()
}

func testSelectWaitNSleep(){
	c := make(chan int)
	o := make(chan bool)
	go func() {
		for {
			select {
			case v := <- c:
				println(v)
			case <- time.After(3 * time.Second):
				println("timeout")
				o <- true
				break
			}
		}
	}()
	<- o
}

func testFlag(){
	propFlag := flag.NewFlagSet("properties", flag.ExitOnError)
	word := propFlag.String("word", "default word", "as string word ")

	isSet := propFlag.Bool("isSet", false, "whether it is set or not?")

	times := propFlag.Int("times", 0 , "The app should run for many times")

	var svar string
	propFlag.StringVar(&svar, "address", ":9990", " server listenning ip address.")



	propFlag.Parse(os.Args[1:])
	fmt.Println("word", *word)
	fmt.Println("isSet", *isSet)
	fmt.Println("address", svar)
	fmt.Println("times: ", *times)
	fmt.Println("tail:", flag.Args())

	if !propFlag.Lookup("isSet").Value.(flag.Getter).Get().(bool) {
		fmt.Println(" `isSet` used the default value 'false' .")
	}


}

func testOsEnv(){
	//os environment
	os.Setenv("registry_ip", "192.168.2.2")
	os.Setenv("default_ip", "127.0.0.1")
	fmt.Println("os env , registry_ip :", os.Getenv("registry_ip"))
	fmt.Println("os env , registry_ip :", os.Getenv("default_ip"))
	//os.Getenv("hold_machine_ip")
	fmt.Println("os env , hold_machine_ip :",len(os.Getenv("hold_machine_ip")) )
	//for _, e := range os.Environ() {
	//	fmt.Println("----env: ", e)
	//}

	cmd , _ := os.Getwd()
	fmt.Println("current workspace rooted directory", cmd)
}

func testUnsafe(){
	fmt.Println(unsafe.Sizeof(float64(0))) // "8"
}

// GetLocalIP returns the non loopback local IP of the host
func getLocalIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println("local ip address: ", ipnet.IP.String())
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func optSig(){
	// Use bitwise OR | to get the bits that are in 1 OR 2
	// 1     = 00000001
	// 2     = 00000010
	// 1 | 2 = 00000011 = 3
	fmt.Println(1 | 2)

	// Use bitwise OR | to get the bits that are in 1 OR 5
	// 1     = 00000001
	// 5     = 00000101
	// 1 | 5 = 00000101 = 5
	fmt.Println(1 | 5)

	// Use bitwise XOR ^ to get the bits that are in 3 OR 6 BUT NOT BOTH
	// 3     = 00000011
	// 6     = 00000110
	// 3 ^ 6 = 00000101 = 5
	fmt.Println(3 ^ 6)

	// Use bitwise AND & to get the bits that are in 3 AND 6
	// 3     = 00000011
	// 6     = 00000110
	// 3 & 6 = 00000010 = 2
	fmt.Println(3 & 6)

	// Use bit clear AND NOT &^ to get the bits that are in 3 AND NOT 6 (order matters)
	// 3      = 00000011
	// 6      = 00000110
	// 3 &^ 6 = 00000001 = 1
	fmt.Println(3 &^ 6)
}

//Limiting Concurrency in Go
func LimitGoroutine(){
	concurrency := 5
	sem := make(chan bool, concurrency)
	urls := []string{"http://jmoiron.net/blog/limiting-concurrency-in-go/", "https://github.com/go-sql-driver/mysql/tree/alloc-pool"}
	for _, url := range urls {
		sem <- true
		go func(url string) {
			defer func() {
				<-sem
			}()
			fmt.Println("imput url:", url )
			// get the url
		}(url)
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
}
//
func main() {
	//fmt.Println("show you name : ", getName("我的年龄", 30, ",你的呢？"))
	//test1()
	//testFibonacci()
	//testFlag()
	//testOsEnv()
	//testUnsafe()
	//getLocalIp()
	//optSig()
	LimitGoroutine()
}
