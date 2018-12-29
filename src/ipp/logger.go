package ipp

import (
	"io/ioutil"
	"log"
)

var Logger = log.New(ioutil.Discard, "", log.Ldate | log.Lmicroseconds | log.Lshortfile)