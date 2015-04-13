package main

import (
	"github.com/akms/suiage/compress"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	var (
		dirPaths             string = "/"
		err                  error
		hostname             string
		beforecheck_fileinfo []os.FileInfo
	)
	if hostname, err = os.Hostname(); err != nil {
		log.Fatal(err)
	}
	hostname = "/mnt/" + hostname
	if _, err = os.Stat(hostname); err != nil {
		if err = os.Mkdir(hostname, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
	if beforecheck_fileinfo, err = ioutil.ReadDir(dirPaths); err != nil {
		log.Fatal(err)
	}
	compress.Compression(beforecheck_fileinfo, dirPaths)
}
