package main

import (
	"crypto/tls"
	"os"
	"testing"

	"toml"
	"go-options"
	"nsqd"
	"fmt"
	"strings"
)

func TestConfigFlagParsing(t *testing.T) {
	opts := nsqd.NewOptions()

	flagSet := nsqdFlagSet(opts)
	flagSet.Parse([]string{})

	var cfg config
	f, err := os.Open("../../contrib/nsqd.cfg")
	if err != nil {
		t.Fatalf("%s", err)
	}
	toml.DecodeReader(f, &cfg)
	cfg.Validate()

	options.Resolve(opts, flagSet, cfg)
	nsqd.New(opts)
	fmt.Println("lookup-address :", strings.Join(opts.NSQLookupdTCPAddresses, ","))

	if opts.TLSMinVersion != tls.VersionTLS10 {
		t.Errorf("min %#v not expected %#v", opts.TLSMinVersion, tls.VersionTLS10)
	}
}
