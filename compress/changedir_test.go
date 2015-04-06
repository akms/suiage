package compress

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestChangeDir(t *testing.T) {
	var dirPath, gopath string
	//GOPATHをとるための悪手
	o, _ := exec.Command(os.Getenv("SHELL"), "-c", "echo $GOPATH").Output()
	gopath = string(o)
	//stringsのTrimRightでchompのような動作
	gopath = strings.TrimRight(gopath, "\n")
	dirPath = gopath + "/src/suiage/compress/test"
	ChangeDir(dirPath)
	workingDir, _ := os.Getwd()
	if workingDir != dirPath {
		t.Errorf("got %s\nwant %s", workingDir, dirPath)
	}
}
