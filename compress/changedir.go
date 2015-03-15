package compress

import (
	"os"
	"log"
)

func ChangeDir(dirName string) {
	err := os.Chdir(dirName)
	if err != nil {
		log.Fatal(err)
	}
}
