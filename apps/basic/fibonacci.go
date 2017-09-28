package main

import (
	"fmt"
	"strconv"
	"log"
	"time"
)

func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

func ExeFibonacci(){
	c := make(chan int)
	quit := make(chan int)
	s := strconv.Quote(`"Fran & Freddie's Diner	☺"`)
	fmt.Println(s)

	//produce data
	go func(){
		for i:= 0; i < 10; i++ {
			log.Println("channel data item ", <- c)
		}
		quit <- 0
	}()
	fibonacci(c, quit)
}


func receiveMsg(){
	c1 := make(chan string)
	c2 := make(chan string)
	t1 := time.Now() // get current time
	go func(){
		//received three message
		for i:=0; i < 3; i++ {
			time.Sleep(time.Millisecond * 150)
			c1 <- "msg 1 with index " + strconv.Itoa(i)
		}
	}()
	go func(){
		for j:=0; j < 3; j++ {
			time.Sleep(time.Millisecond * 200)
			c2 <- "msg 2 with index "+ strconv.Itoa(j)
		}
	}()

	//print two message
	for i:= 0; i < 3; i++ {
		select {
		case msg1 := <- c1:
			log.Println("received msg :", msg1)
		case msg2 := <- c2:
			log.Println("received msg :", msg2)
		}
	}
	elapsed := time.Since(t1)
	log.Println("time elapsed ", elapsed)

}
