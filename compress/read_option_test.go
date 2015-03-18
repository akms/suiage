package compress

import (
	"os"
	"testing"
)

func TestReadOption(t *testing.T) {
	var (
		read_strings []string
		test_strings = []string{"^_old$", "^boot$", "^opt$", "^root$", "^sbin$", "^etc$", "^var$", "^home$"}
	)
	read_strings = ReadOption()
	workingDir, _ := os.Getwd()
	if workingDir != "/etc" {
		t.Errorf("workingdir is not /etc. now workingdir is %s", workingDir)
	}
	for i, s := range read_strings {
		if s != test_strings[i] {
			t.Errorf("test fatal got %s", s)
		}
	}

}
