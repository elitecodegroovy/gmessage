package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"time"
	"util"
)

//define the small file size
const maxSizeOfSmallFile = 1024 * 512

//allow max goroutine
const maxGoroutines = 50000

type pathsInfo struct {
	size  int64
	paths []string
}

type fileInfo struct {
	sha1 []byte
	size int64
	path string
}

func searchDupFile(infoChan chan fileInfo, dirname string) {
	var runSync util.WaitGroupWrapper
	filepath.Walk(dirname, doDeeplyWalkFunc(infoChan, &runSync))
	runSync.Wait() // Blocks until all the work is done
	close(infoChan)
}

func doDeeplyWalkFunc(infoChan chan fileInfo, runSync *util.WaitGroupWrapper) func(string, os.FileInfo, error) error {
	return func(path string, info os.FileInfo, err error) error {
		if err == nil && info.Size() > 0 && (info.Mode()&os.ModeType == 0) {
			//It defines the limit to the usage of the device's resources.
			if info.Size() < maxSizeOfSmallFile ||
				runtime.NumGoroutine() > maxGoroutines {
				checkFile(path, info, infoChan, nil)
			} else {
				runSync.Wrap(func() {
					checkFile(path, info, infoChan, nil)
				})
			}
		}
		return nil // We ignore all errors
	}
}

func checkFile(filename string, info os.FileInfo, infoChan chan fileInfo, done func()) {
	if done != nil {
		defer done()
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Println("error:", err)
		return
	}
	defer file.Close()
	hash := sha256.New()
	if size, err := io.Copy(hash, file); size != info.Size() || err != nil {
		if err != nil {
			log.Println("error:", err)
		} else {
			log.Println("error: failed to read the whole file:", filename)
		}
		return
	}
	infoChan <- fileInfo{hash.Sum(nil), info.Size(), filename}
}

func mergeFileInfos(infoChan <-chan fileInfo) map[string]*pathsInfo {
	pathData := make(map[string]*pathsInfo)
	format := fmt.Sprintf("%%016X:%%%dX", sha256.Size*2) // == "%016X:%40X"
	for info := range infoChan {
		key := fmt.Sprintf(format, info.size, info.sha1)
		value, found := pathData[key]
		if !found {
			value = &pathsInfo{size: info.size}
			pathData[key] = value
		}
		value.paths = append(value.paths, info.path)
	}
	return pathData
}

func outputDupFileInfos(pathData map[string]*pathsInfo) {
	keys := make([]string, 0, len(pathData))
	for key := range pathData {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		value := pathData[key]
		if len(value.paths) > 1 {
			fmt.Printf("%d duplicatedFiles (%s bytes):\n",
				len(value.paths), commas(value.size))
			sort.Strings(value.paths)
			for _, name := range value.paths {
				fmt.Printf("\t%s\n", name)
			}
		}
	}
}

// commas() returns a string representing the whole number with comma
// grouping.
func commas(x int64) string {
	value := fmt.Sprint(x)
	for i := len(value) - 3; i > 0; i -= 3 {
		value = value[:i] + "," + value[i:]
	}
	return value
}

func doStartSearchDupFiles() {
	t1 := time.Now()
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		fmt.Printf("usage: %s <path>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	infoChan := make(chan fileInfo, maxGoroutines*2)
	if len(os.Args) == 1 {
		go searchDupFile(infoChan, "./")
	} else {
		go searchDupFile(infoChan, os.Args[1])
	}

	pathData := mergeFileInfos(infoChan)
	outputDupFileInfos(pathData)
	fmt.Printf("time elapses %d ms", time.Since(t1).Nanoseconds()/1000000)
}

//func main() {
//	doStartSearchDupFiles()
//}
