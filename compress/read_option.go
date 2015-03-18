package compress

import (
	"bufio"
	"log"
	"os"
)

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func ReadOption() (lines []string) {
	ChangeDir("/etc")
	if FileExists("suiage.conf") {
		infile_options, err := os.Open("suiage.conf")
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
		return
	} else {
		lines = []string{""}
		return
	}
}
