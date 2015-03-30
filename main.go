package main

import (
	"log"
	"os"
	"suiage/compress"
	//"github.com/akms/suiage/compress"
)

func main() {
	var (
		dirPaths string = "/"
		err      error
		hostname string
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
	compress.CheckTarget(dirPaths)
}
