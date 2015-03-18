package compress

import (
	"log"
	"os"
)

func ChangeDir(dirName string) {
	err := os.Chdir(dirName)
	if err != nil {
		log.Fatal(err)
	}
}
