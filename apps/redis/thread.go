package main


import (
	"sync"
	"runtime"
	"time"
	"io/ioutil"
	"os"
	"fmt"
)

type JobFunc func(int, interface{}, chan interface{})

func threadMain(id int, queue chan interface {}, wg *sync.WaitGroup, job JobFunc) chan bool {
	quitCommand := make(chan bool, 1)
	go func() {
		for {
			select {
			case task := <-queue:
				wg.Add(1)
				job(id, task, queue)
				wg.Done()
			case <-quitCommand:
				return
			}

		}
	}()
	return quitCommand
}

func Concurrent(queue chan interface {}, job JobFunc) {
	var wg sync.WaitGroup
	cpuCount := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuCount)

	quitCommands := make([]chan bool, cpuCount)
	for i := 0; i < cpuCount; i++ {
		quitCommands[i] = threadMain(i + 1, queue, &wg, job)
	}

	ticker := time.Tick(time.Millisecond * 10)

	for _ = range ticker {
		if len(queue) == 0 {
			for _, quitCommand := range quitCommands {
				quitCommand <- true
			}
			wg.Wait()
			break
		}
	}
}

//
//func main() {
//	fileInfos, _ := ioutil.ReadDir(os.Args[1])
//	queue := make(chan interface{}, len(fileInfos))
//
//	for _, fileInfo := range fileInfos {
//		queue <- fileInfo
//	}
//
//	Concurrent(queue, func (id int, task interface {}, queue chan interface {}) {
//		fileInfo := task.(os.FileInfo)
//		fmt.Printf(">> thread %d: %s\n", id, fileInfo.Name())
//		// do some work to use fileInfo
//		fmt.Printf("<< thread %d\n", id)
//	})
//}
