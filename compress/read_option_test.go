package compress

import (
	"regexp"
	"testing"
)

func TestReadOption(t *testing.T) {
	var (
		check_strings    = []string{"^_old$", "^usr$", "^lib$", "^etc$", "^data$", "^opt$"}
		read_strings     []string
		fchecker         bool
		gopath, fullpath string
	)

	gopath = getGopath()
	fullpath = gopath + "/src/suiage/compress/test/suiage.conf"

	if _, err := ReadOption("/etc/sugiage.conf"); err == nil {
		t.Errorf("check filename faild")
	}

	read_strings, _ = ReadOption(fullpath)
	comment_Regexp := regexp.MustCompile(`^#`)
	nilstr_Regexp := regexp.MustCompile(`^$`)
	for _, c := range check_strings {
		fchecker = false
		for _, r := range read_strings {
			if comment_Regexp.MatchString(r) || nilstr_Regexp.MatchString(r) {
				t.Errorf("catch ng string %s or blank", r)
			}
			if c == r {
				fchecker = true
			}
		}
		if !fchecker {
			t.Errorf("can't find word %s", c)
		}
	}
}
