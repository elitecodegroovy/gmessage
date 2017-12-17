package main

import (
	"bytes"
	"strconv"
	"time"
	"fmt"
)

func buffString(){
	t1 := time.Now().Local()
	var buffer bytes.Buffer
	for i:=0; i < 10000; i++ {
		buffer.WriteString(""+ strconv.Itoa(i))
	}
	fmt.Println("bytes.Buffer elapses: ", time.Since(t1))

	t2 := time.Now().Local()
	str := ""
	for j:=0; j < 10000; j++ {
		str += strconv.Itoa(j)
	}
	fmt.Println("String concatenate elapses: ", time.Since(t2))
}

//func main(){
//	buffString()
//}
