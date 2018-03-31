package main

import (
	"os"
	"encoding/csv"
	"github.com/elitecodegroovy/gmessage/apps/osext"
	"bufio"
	"io"
	"fmt"
	"log"
	"flag"
)
func initWorkingDirectory() string {
	var customPath string
	// Check if a custom path has been provided by the user.
	flag.StringVar(&customPath, "custom-path", "",
		"Specify a custom path to the asset files. This needs to be an absolute path.")
	flag.Parse()
	// Get the absolute path this executable is located in.
	executablePath, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal("Error: Couldn't determine working directory: " + err.Error())
	}
	// Set the working directory to the path the executable is located in.
	os.Chdir(executablePath)
	// Return the user-specified path. Empty string if no path was provided.
	return customPath
}

func readAndOutput(){
	csvFile, _ := os.Open( "file.csv")
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	fmt.Println("[")
	i := 0
	for {
		line, err := reader.Read()
		i++
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal("read error " + err.Error())
		}
		if i == 1 {
			fmt.Printf("[%s, %s,8,1]", line[0], line[1])
		}else {
			fmt.Printf(",\n[%s, %s,8,1]", line[0], line[1])
		}
	}
	fmt.Println("[")
}

func main(){
	readAndOutput()
}
