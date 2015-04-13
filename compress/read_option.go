package compress

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func FileExists(filename string) bool {
	var err error
	_, err = os.Stat(filename)
	if err != nil {
		return false
	}
	return true
}

func ReadOption(fullpath string) (lines []string, err error) {
	var (
		dirpath, filename string
		infile_options    *os.File
		fi                os.FileInfo
		size, n           int64
		serr              error
	)
	dirpath, filename = filepath.Split(fullpath)
	ChangeDir(dirpath)
	if filename != "suiage.conf" {
		err = fmt.Errorf("file name is not suiage.conf got %s", filename)
		lines = make([]string, 0, 0)
		return
	}
	if FileExists(filename) {
		infile_options, err = os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer infile_options.Close()
		if fi, err = infile_options.Stat(); err == nil {
			if size = fi.Size(); size < 1e9 {
				n = size + bytes.MinRead
			}
		}
		lines = make([]string, 0, n)
		comment_Regexp := regexp.MustCompile(`^#`)
		nilstr_Regexp := regexp.MustCompile(`^$`)
		scanner := bufio.NewScanner(infile_options)
		for scanner.Scan() {
			if !comment_Regexp.MatchString(scanner.Text()) && !nilstr_Regexp.MatchString(scanner.Text()) {
				lines = append(lines, scanner.Text())
			}
		}
		if serr = scanner.Err(); serr != nil {
			log.Fatal(serr)
		}
		return lines, err
	} else {
		lines = make([]string, 0, 0)
		return
	}
}
