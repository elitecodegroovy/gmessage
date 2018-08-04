package main

import (
	"flag"
	"log"

	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/elitecodegroovy/gmessage/gio"
	"github.com/kardianos/osext"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	// Initialization of the working directory. Needed to load asset files.
	binaryFilePath = initWorkingDirectory()
)

// NOTE: Use tls scheme for TLS, e.g. nats-pub -s tls://demo.nats.io:4443 foo hello
func usage() {
	log.Fatalf("Usage: nats-pub [-s server (%s)] <subject> <msg> \n", nats.DefaultURL)
}

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

func doPublish(nc *nats.Conn, subj string) error {
	filename := filepath.Join(binaryFilePath, "file.csv")
	fmt.Println("[]:" + filename)
	csvFile, err := os.Open(filename)
	if err != nil {
		fmt.Errorf("openfile error in func doPublish %s", err.Error())
		return err
	}
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	i := 0
	for {
		line, err := reader.Read()
		i++
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal("read error " + err.Error())
			return err
		}
		fmt.Println("msg:" + strings.Join(line, "") + "\t" + string(i))
		nc.Publish(subj, []byte(strings.Join(line, "")))
	}
	return nil
}

func publishMsg(subj string) {
	var urls = flag.String("s", "nats://192.168.1.225:6222", "The nats server URLs (separated by comma)")

	log.SetFlags(0)
	//flag.Usage = usage
	//flag.Parse()

	//args := flag.Args()
	//if len(args) < 2 {
	//	usage()
	//}

	nc, err := nats.Connect(*urls)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	doPublish(nc, subj)
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Published [%s] : \n", subj)
	}
}

func main() {
	//subj : topic message
	publishMsg("test01")
}
