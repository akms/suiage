package compress

import (
	"bufio"
	"os"
	"log"
)

func ReadOption() []string {
	ChangeDir("/etc")
	infile_options,err := os.Open("suiage.conf")
	if err != nil {
		log.Fatal(err)
	}
	defer infile_options.Close()
	lines := make([]string, 0, 100)
	scanner := bufio.NewScanner(infile_options)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if serr := scanner.Err(); serr != nil {
		log.Fatal(serr)
	}
	return lines
}
