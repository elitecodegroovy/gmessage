package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"time"
)

var cores = runtime.NumCPU()

type LineInfo struct {
	filename string
	lino     int
	line     string
}

type Task struct {
	filename  string
	LineInfos chan<- LineInfo
}

func (task Task) Do(lineRx *regexp.Regexp) {
	file, err := os.Open(task.filename)
	if err != nil {
		log.Printf("error: %s\n", err)
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for lino := 1; ; lino++ {
		line, err := reader.ReadBytes('\n')
		line = bytes.TrimRight(line, "\n\r")
		if lineRx.Match(line) {
			task.LineInfos <- LineInfo{task.filename, lino, string(line)}
		}
		if err != nil {
			if err != io.EOF {
				log.Printf("error:%d: %s\n", lino, err)
			}
			break
		}
	}
}

func scheduleTasks(done chan<- struct{}, lineRx *regexp.Regexp, tasks <-chan Task) {
	for i := 0; i < cores; i++ {
		go doTasks(done, lineRx, tasks)
	}
}

func doGrep(timeout int64, lineRx *regexp.Regexp, filenames []string) {
	tasks := make(chan Task, cores)
	LineInfos := make(chan LineInfo, min(1000, len(filenames)))
	done := make(chan struct{}, cores)

	go addTasks(tasks, filenames, LineInfos)
	scheduleTasks(done, lineRx, tasks)
	waitAndProcessResults(timeout, done, LineInfos)
}

func addTasks(tasks chan<- Task, filenames []string, LineInfos chan<- LineInfo) {
	for _, filename := range filenames {
		tasks <- Task{filename, LineInfos}
	}
	close(tasks)
}

func doTasks(done chan<- struct{}, lineRx *regexp.Regexp, tasks <-chan Task) {
	for task := range tasks {
		task.Do(lineRx)
	}
	done <- struct{}{}
}

func waitAndProcessResults(timeout int64, done <-chan struct{}, LineInfos <-chan LineInfo) {
	finish := time.After(time.Duration(timeout))
	for working := cores; working > 0; {
		select { // Blocking
		case LineInfo := <-LineInfos:
			fmt.Printf("%s:%d:%s\n", LineInfo.filename, LineInfo.lino, LineInfo.line)
		case <-finish:
			fmt.Println("timed out")
			return // Time's up so finish with what LineInfos there were
		case <-done:
			working--
		}
	}
	for {
		select { // Nonblocking
		case LineInfo := <-LineInfos:
			fmt.Printf("%s:%d:%s\n", LineInfo.filename, LineInfo.lino,
				LineInfo.line)
		case <-finish:
			fmt.Println("timed out")
			return // Time's up so finish with what LineInfos there were
		default:
			return
		}
	}
}

func min(x int, ys ...int) int {
	for _, y := range ys {
		if y < x {
			x = y
		}
	}
	return x
}

func startChannelPattern() {
	var timeoutOpt *int64 = flag.Int64("timeout", 0, "seconds (0 means no timeout)")
	flag.Parse()
	if *timeoutOpt < 0 || *timeoutOpt > 240 {
		log.Fatalln("timeout must be in the range [0,240] seconds")
	}
	if len(os.Args) < 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("usage: %s <regexp> <files>, e.g. channel.exe main .\\ \n",
			filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	var timeout int64 = 1e9 * 60 * 10 // 10 minutes!
	if *timeoutOpt != 0 {
		timeout = *timeoutOpt * 1e9
	}
	if lineRx, err := regexp.Compile(os.Args[1]); err != nil {
		log.Fatalf("invalid regexp: %s\n", err)
	} else if len(os.Args) == 2 {
		//current file directory files
		doGrep(timeout, lineRx, ExecSearch())
	} else {
		doGrep(timeout, lineRx, commandLineFiles(os.Args[2:]))
	}
}

//func main() {
//	startChannelPattern()
//}
