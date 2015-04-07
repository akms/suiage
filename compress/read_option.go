package compress

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return true
}

func ReadOption(fullpath string) (lines []string, err error) {
	var (
		dirpath, filename string
	)
	dirpath, filename = filepath.Split(fullpath)
	ChangeDir(dirpath)
	if filename != "suiage.conf" {
		err = fmt.Errorf("file name is not suiage.conf got %s", filename)
		lines = make([]string, 0, 0)
		return
	}
	if FileExists(filename) {
		infile_options, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer infile_options.Close()
		lines = make([]string, 0, 100)
		scanner := bufio.NewScanner(infile_options)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if serr := scanner.Err(); serr != nil {
			log.Fatal(serr)
		}
		return lines,err
	} else {
		lines = make([]string, 0, 0)
		return
	}
}
