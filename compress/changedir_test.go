package compress

import (
	"os"
	"testing"
)

func TestChangeDir(t *testing.T) {
	var dirPath, gopath string

	gopath = getGopath()
	dirPath = gopath + "/src/suiage/compress/test"
	ChangeDir(dirPath)
	workingDir, _ := os.Getwd()
	if workingDir != dirPath {
		t.Errorf("got %s\nwant %s", workingDir, dirPath)
	}
}
