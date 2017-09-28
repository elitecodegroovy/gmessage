package main

import (
	"internal/app"
	"flag"
	"os"
	"log"
	"fmt"
)

var (
	flagSet = flag.NewFlagSet("input_array", flag.ExitOnError)
	xArray = app.StringArray{}
)

func init(){
	flagSet.Var(&xArray, "x-array", "input string array with comma separator(may be given multiple times)")
}

func main(){
	flagSet.Parse(os.Args[1:])
	if len(xArray) == 0 {
		log.Fatal("x-array CL input parameter must be provided. \n e.g. --x-array=1 --x-array=2")
	}
	log.Output(2, fmt.Sprintf("lookupdHTTPAddrs length: %d, %v", len(xArray), xArray))
}
