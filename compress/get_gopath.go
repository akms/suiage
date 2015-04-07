package compress

import (
	"os"
	"os/exec"
	"strings"
)


func getGopath() (gopath string) {
	//GOPATHをとるための悪手
	o, _ := exec.Command(os.Getenv("SHELL"), "-c", "echo $GOPATH").Output()
	gopath = string(o)
	//stringsのTrimRightでchompのような動作
	gopath = strings.TrimRight(gopath, "\n")
	return 
}
