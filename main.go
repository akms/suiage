package main

import (
	"log"
	"os"
	"suiage/compress"
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
	if err = os.Mkdir(hostname, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	compress.CheckTarget(dirPaths)
}
