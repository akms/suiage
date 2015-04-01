package compress

import (
	"fmt"
	"regexp"
	"strings"
)

type Matcher interface {
	MatchDefaultTarget() bool
	MatchOptionTarget() bool
}

type Target struct {
	name string
}

var (
	default_except_targets, option_except_targets []string = strings.Fields(`^lost\+found$ ^proc$ ^sys$ ^dev$ ^mnt$ ^media$ ^run$ ^selinux$ ^boot$ ^_old$`), ReadOption()
)

func (target *Target) MatchDefaultTarget() bool {
	for i, s := range default_except_targets {
		default_Regexp := regexp.MustCompile(s)
		if default_Regexp.MatchString(target.name) {
			default_except_targets = append(default_except_targets[:i], default_except_targets[i+1:]...)
			return false
		}
	}
	return true
}

func (target *Target) MatchOptionTarget() bool {
	for _, s := range option_except_targets {
		option_Regexp := regexp.MustCompile(s)
		if option_Regexp.MatchString(target.name) {
			return true
		}
	}
	return false
}

func targetMatch(matcher Matcher) bool {
	if !matcher.MatchDefaultTarget() {
		fmt.Println("Target match default")
		return true
	}
	if matcher.MatchOptionTarget() {
		fmt.Println("Target match option")
		return true
	}
	return false
}
