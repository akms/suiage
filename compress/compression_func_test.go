package compress

import (
	"testing"
	"archive/tar"
	"compress/gzip"
	"os"
	"strings"
)

func TestMakeFile(t *testing.T) {
	var (
		gw *gzip.Writer
		tw *tar.Writer
		file *os.File
	)
	gw, tw, file = MakeFile()
	if gw == nil {
		t.Errorf("make faild gzip writer")
	}
	if tw == nil {
		t.Errorf("make faild tar writer")
	}
	if file == nil {
		t.Errorf("make faild file")
	}
}

func TestMatchDefaultTarget(t *testing.T) {
	str := strings.Fields(`^lost\+found$ ^proc$ ^sys$ ^dev$ ^mnt$ ^media$ ^run$ ^selinux$`)
	for _, s := range str {
		if !MatchDefaultTarget(s) {
			t.Errorf("Match faild %s",s)
		}
	}
}


func TestCheckTarget(t *testing.T) {
	
}
