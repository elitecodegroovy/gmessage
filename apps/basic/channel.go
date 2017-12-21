package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"util"
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
		if strbry_closed && choco_closed {
			return
		}
		fmt.Println("Waiting for a new cake ...")
		select {
		case cakeName, strbry_ok := <-strbry_cs:
			if !strbry_ok {
				strbry_closed = true
				fmt.Println(" ... Strawberry channel closed!")
			} else {
				fmt.Println("Received from Strawberry channel.  Now packing", cakeName)
			}
		case cakeName, choco_ok := <-choco_cs:
			if !choco_ok {
				choco_closed = true
				fmt.Println(" ... Chocolate channel closed!")
			} else {
				fmt.Println("Received from Chocolate channel.  Now packing", cakeName)
			}
		}
	}
}

func sendReceiveChan() {
	strawBerryChan := make(chan string)
	chocolateChan := make(chan string)

	//two cake makers
	var syncWrapper util.WaitGroupWrapper
	syncWrapper.Wrap(func() {
		makeCakeAndSend(strawBerryChan, "Chocolate", 3) //make 3 chocolate cakes and send
	})
	syncWrapper.Wrap(func() {
		makeCakeAndSend(chocolateChan, "Strawberry", 3) //make 3 strawberry cakes and send
	})

	//one cake receiver and packer
	syncWrapper.Wrap(func() {
		receiveCakeAndPack(strawBerryChan, chocolateChan) //pack all cakes received on these cake channels
	})
	syncWrapper.Wait()
	fmt.Println("end select and channel ....")
}

const (
	c1 = imag(2i)                   // imag(2i) = 2.0 is a constant
	c2 = len([10]float64{2})        // [10]float64{2} contains no function calls
	c3 = len([10]float64{c1})       // [10]float64{c1} contains no function calls
	c4 = len([10]float64{imag(2i)}) // imag(2i) is a constant and no function call is issued
	//c5 = len([10]float64{imag(z)})   // invalid: imag(z) is a (non-constant) function call
)

func testPathJoin() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	gopath := filepath.Clean(filepath.Join(cwd, "../../../../"))
	log.Println("GOPATH is", gopath)

	log.Println("c2", c2, ", c3:", c3, "c4: ", c4)

}

//func main() {
//	//sendReceiveChan
//	//ExeFibonacci()
//	//receiveMsg()
//	testPathJoin()
//}
