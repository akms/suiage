package compress

import (
	"log"
	"os"
)

func ChangeDir(dirName string) {
	var err error
	err = os.Chdir(dirName)
	if err != nil {
		log.Fatal(err)
	}
}
