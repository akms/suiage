package compress

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"testing"
)

func MakeCopy() {
	var (
		buf                    bytes.Buffer
		copy_fileWriter        io.WriteCloser
		err                    error
		origin_file, copy_file *os.File
	)
	origin_file, err = os.Open("/etc/suiage.conf")
	if err != nil {
		log.Fatal("can't open file")
	}
	io.Copy(&buf, origin_file)
	defer origin_file.Close()
	copy_file, err = os.Create("/tmp/t_suiage.conf")
	if err != nil {
		log.Fatal(err)
	}
	copy_fileWriter = copy_file
	defer copy_file.Close()
	copy_fileWriter.Write(buf.Bytes())
	defer copy_fileWriter.Close()
}

func TestReadOption(t *testing.T) {
	var (
		read_strings []string
		fchecker     bool
	)

	MakeCopy()
	read_strings = ReadOption()
	f, _ := os.Open("/tmp/t_suiage.conf")
	defer f.Close()
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		s := scan.Text()
		fchecker = false
		for _, r := range read_strings {
			if s == r {
				fchecker = true
			}
		}
		if !fchecker {
			t.Errorf("can't find word %s", s)
		}
	}
	os.Remove("/tmp/t_suiage.conf")
	workingDir, _ := os.Getwd()
	if workingDir != "/etc" {
		t.Errorf("workingdir is not /etc. now workingdir is %s", workingDir)
	}

}
