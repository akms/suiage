package compress

import (
	"bytes"
	"io"
	"os"
	"testing"
	"log"
)

func MakeCopy() {
	var (
		b           bytes.Buffer
		cfileWriter io.WriteCloser
	)
	origin, err := os.Open("/etc/suiage.conf")
	if err != nil {
		log.Fatal("can't open file")
	}
	io.Copy(&b, origin)
	defer origin.Close()
	cfile, _ := os.Create("/tmp/t_suiage.conf")
	cfileWriter = cfile
	defer cfile.Close()
	cfileWriter.Write(b.Bytes())
	defer cfileWriter.Close()
}

func TestReadOption(t *testing.T) {
	var (
		read_strings []string
		test_strings []string
	)


	read_strings = ReadOption()

	workingDir, _ := os.Getwd()
	if workingDir != "/etc" {
		t.Errorf("workingdir is not /etc. now workingdir is %s", workingDir)
	}
	
	
	
	
	
	

}
