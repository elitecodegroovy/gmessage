package main

import (
	"os"
	"flag"
	"fmt"
	"path/filepath"
	"sync"
)

func ReadDirNumSort(dirname string, reverse bool) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	return list, nil
}


func getCmdDirName()(dirname string){
	flag.StringVar(&dirname, "dirname", "./", "maximum file size (-1 means no maximum)")
	flag.Parse()
	return dirname
}

func searchFile(searchDir string)([]string){
	fileNamePaths := []string{}
	err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			fileNamePaths = append(fileNamePaths, path)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
	for _, fileName := range fileNamePaths {
		fmt.Println(fileName)
	}
	return fileNamePaths
}

type WaitGroupWrapper struct {
	sync.WaitGroup
}

func (w *WaitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}

func ExecSearch() []string{
	dirname := getCmdDirName()
	fmt.Println("input dirname:", dirname)
	absPath, err:= filepath.Abs(dirname)
	if err != nil {
		panic(err)
	}
	return searchFile(absPath)
}

func execConcurrentSearch(){
	var f *os.File
	var t WaitGroupWrapper
	fileNames := []interface{}{}
	fileNameChan := make(chan string)
	dirname := getCmdDirName()
	fmt.Println("input dirname:", dirname)
	absPath, err:= filepath.Abs(dirname)
	if err != nil {
		panic(err)
	}
	f, err = os.Open(absPath)
	if err != nil {
		panic(err)
	}
	list, err := f.Readdir(-1)
	defer f.Close()
	for _, fPath := range list {
		if fPath.IsDir() {
			t.Wrap(func(){
				fileNames = append(fileNames, searchFile(fPath.Name()))
			})
		}else {
			fileNameChan <- fPath.Name()
			fmt.Println(fPath.Name())
		}
	}
	stop := false
	for{
		if stop {
			break
		}
		select {
		case  o, e := <- fileNameChan :
			if !e {
				stop = true
			}
			fmt.Println(o)
		}
	}
	t.Wait()
	fmt.Println("all file has been search!")
}

//func main(){
//	//TODO ...fixed the issue.
//	//execConcurrentSearch()
//	execSearch()
//}