package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/kardianos/osext"
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	// Initialization of the working directory. Needed to load asset files.
	binaryFilePath = initWorkingDirectory()
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

func readAndOutput() {
	filename := filepath.Join(binaryFilePath, "file.csv")
	fmt.Println("[]:" + filename)
	csvFile, _ := os.Open(filename)
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	fmt.Println("[")
	i := 0
	for {
		line, err := reader.Read()
		i++
		if err == io.EOF {
			fmt.Println("]")
			break
		} else if err != nil {
			log.Fatal("read error " + err.Error())
		}
		if i == 1 {
			fmt.Printf("[%s,%s,8,1]", line[0], line[1])
		} else {
			fmt.Printf(",\n[%s,%s,8,1]", line[0], line[1])
		}
	}

}

func main() {
	readAndOutput()
}
