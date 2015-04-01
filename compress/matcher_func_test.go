package compress

import (
	"strings"
	"testing"
)

func TestMatchDefaultTarget(t *testing.T) {
	var target *Target = &Target{}
	default_except_targets = strings.Fields(`^lost\+found$ ^proc$ ^sys$ ^dev$ ^mnt$ ^media$ ^run$ ^selinux$`)
	str := strings.Fields("lost+found proc sys dev mnt media selinux run")
	for _, s := range str {
		target.name = s
		if target.MatchDefaultTarget() {
			t.Errorf("Match faild %s", target.name)
		}
	}
	default_except_targets = strings.Fields(`^lost\+found$ ^proc$ ^sys$ ^dev$ ^mnt$ ^media$ ^run$ ^selinux$ ^run$`)
	str = strings.Fields("var etc lib lib64")
	for _, s := range str {
		target.name = s
		if !target.MatchDefaultTarget() {
			t.Errorf("Match faild %s", target.name)
		}
	}
}

func TestMatchOptionTarget(t *testing.T) {
	var target *Target = &Target{}
	str := strings.Fields("_old etc var src")
	option_except_targets = strings.Fields("^_old$ ^etc$ ^var$ ^src$")
	for _, s := range str {
		target.name = s
		if !target.MatchOptionTarget() {
			t.Errorf("Match faild %s.", target.name)
		}
	}
	str = strings.Fields("lost+found proc sys dev mnt media selinux run")
	for _, s := range str {
		target.name = s
		if target.MatchOptionTarget() {
			t.Errorf("Match faild %s", target.name)
		}
	}
}

func TestTargetMatch(t *testing.T) {
	var target *Target = &Target{"_old"}
	option_except_targets = strings.Fields("^_old$ ^etc$")
	if !targetMatch(target) {
		t.Errorf("Faild")
	}
	var fileio *Fileio = &Fileio{Target:target}
	if !targetMatch(fileio) {
		t.Errorf("Faild")
	}
	var matcher Matcher = &Target{"etc"}
	if !targetMatch(matcher) {
		t.Errorf("Faild")
	}
}
