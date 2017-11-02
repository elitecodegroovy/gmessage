package main

import (
	"fmt"
	"strconv"
	"util"
	"os"
	"path/filepath"
	"log"
)

func makeCakeAndSend(cs chan string, flavor string, count int) {
	for i := 1; i <= count; i++ {
		cakeName := flavor + " Cake " + strconv.Itoa(i)
		cs <- cakeName //send a strawberry cake
	}
	close(cs)
}

func receiveCakeAndPack(strbry_cs chan string, choco_cs chan string) {
	strbry_closed, choco_closed := false, false

	for {
		//if both channels are closed then we can stop
		if (strbry_closed && choco_closed) { return }
		fmt.Println("Waiting for a new cake ...")
		select {
		case cakeName, strbry_ok := <-strbry_cs:
			if (!strbry_ok) {
				strbry_closed = true
				fmt.Println(" ... Strawberry channel closed!")
			} else {
				fmt.Println("Received from Strawberry channel.  Now packing", cakeName)
			}
		case cakeName, choco_ok := <-choco_cs:
			if (!choco_ok) {
				choco_closed = true
				fmt.Println(" ... Chocolate channel closed!")
			} else {
				fmt.Println("Received from Chocolate channel.  Now packing", cakeName)
			}
		}
	}
}

func sendReceiveChan(){
	strawBerryChan := make(chan string)
	chocolateChan := make(chan string)

	//two cake makers
	var syncWrapper util.WaitGroupWrapper
	syncWrapper.Wrap(func() {
		makeCakeAndSend(strawBerryChan, "Chocolate", 3)  //make 3 chocolate cakes and send
	})
	syncWrapper.Wrap(func(){
		makeCakeAndSend(chocolateChan, "Strawberry", 3)  //make 3 strawberry cakes and send
	})

	//one cake receiver and packer
	syncWrapper.Wrap(func(){
		receiveCakeAndPack(strawBerryChan, chocolateChan)  //pack all cakes received on these cake channels
	})
	syncWrapper.Wait()
	fmt.Println("end select and channel ....")
}

func testPathJoin(){
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	gopath := filepath.Clean(filepath.Join(cwd, "../../../../"))
	log.Println("GOPATH is", gopath)
}
func main() {
	//sendReceiveChan
	//ExeFibonacci()
	//receiveMsg()
	testPathJoin()
}