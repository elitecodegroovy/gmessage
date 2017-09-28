package main

import (
	"flag"
	"go-options"
	"os"
	"os/signal"
	"syscall"
	"nsqplus"
	"log"
	"toml"
	"internal/app"
	"fmt"
	"internal/version"
)

var (
	flagSet 		= flag.NewFlagSet("nsqplus", flag.ExitOnError)
	config      		= flagSet.String("config", "", "path to config file")

	showVersion 		= flagSet.Bool("version", false, "print version string")
	httpAddress      	= flagSet.String("http-address", "0.0.0.0:8800", "<addr>:<port> to listen on for TCP clients")

	nsqdHTTPAddresses       = app.StringArray{}
)

func init() {
	flagSet.Var(&nsqdHTTPAddresses, "nsqd-http-address", "nsqd HTTP address (may be given multiple times)")
}
func main() {
	flagSet.Parse(os.Args[1:])

	if *showVersion {
		fmt.Println(version.String("nsq-plus"))
		return
	}

	exitChan := make(chan int)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		exitChan <- 1
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	//logic entrance
	var cfg map[string]interface{}
	if *config != "" {
		_, err := toml.DecodeFile(*config, &cfg)
		if err != nil {
			log.Fatalf("ERROR: failed to load config file %s - %s", *config, err)
		}
	}

	opts := nsqplus.NewOptions()
	options.Resolve(opts, flagSet, cfg)

	nsqplusMain := nsqplus.New(opts)

	nsqplusMain.Main()
	<-exitChan
	nsqplusMain.Exit()
}
