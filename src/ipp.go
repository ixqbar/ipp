package main

import (
	"flag"
	"fmt"
	"ipp"
	"os"
	"strings"
)

var version = flag.Bool("version", false, "print current version")
var showLocalIp = flag.Bool("ip", false, "print current machine ip")
var runAsServer = flag.Bool("server", false, "run as server mode")
var runAsClient = flag.Bool("client", false, "run as client mode")
var verbose = flag.Bool("verbose", false, "show running logs")
var help = flag.Bool("help", true, "show help")

func usage() {
	fmt.Printf("Usage: %s options\nOptions:\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *verbose {
		ipp.Logger.SetOutput(os.Stdout)
	}

	if *version {
		fmt.Printf("%s\n", ipp.VERSION)
		os.Exit(0)
	} else if *showLocalIp {
		fmt.Printf("current machine all ips:\n%s\n", strings.Join(ipp.GetCurrentMachineIps(), "\n"))
		os.Exit(0)
	} else if *runAsClient {
		ipp.Run(false)
	} else if *runAsServer {
		ipp.Run(true)
	} else if *help {
		usage()
	}
}
