package channel

import (
	"fmt"
	"strconv"
	"time"
)

func addFunc(fxChan chan func() string) {
	i := 0
	max := 10
	for {
		i++
		fxChan <- func() string {
			return "msg " + strconv.Itoa(i)
		}
		time.Sleep(1 * time.Second)
		if i > max {
			fmt.Println("close addFunc()")
			close(fxChan)
			return
		}
	}

}

//It is always doing.
func doTask(fxChan chan func() string) {
	i := 10
	for {
		select {
		case rfx := <-fxChan:
			msg := rfx()
			fmt.Println(msg, "Received!")
		default:
			i++
			fmt.Println(".")
		}
		if i == 10 {
			return
		}

	}
}

func DoFuncChan() {
	fxChan := make(chan func() string)
	go addFunc(fxChan)
	doTask(fxChan)
}
