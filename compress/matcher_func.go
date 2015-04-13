package compress

import (
	"regexp"
	"strings"
	"log"
)

type Matcher interface {
	MatchDefaultTarget() bool
	MatchOptionTarget() bool
	setMatcherName(string)
}

type Target struct {
	name string
}

var (
	default_except_targets []string  = strings.Fields(`^lost\+found$ ^proc$ ^sys$ ^dev$ ^mnt$ ^media$ ^run$ ^selinux$ ^_old$`)
	option_except_targets,err = ReadOption("/etc/suiage.conf")
)

func (target *Target) setMatcherName(s string) {
	target.name = s
}

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
	if err != nil {
		log.Fatal(err)
	}
	for _, s := range option_except_targets {
		option_Regexp := regexp.MustCompile(s)
		if option_Regexp.MatchString(target.name) {
			return true
		}
	}
	return false
}

func SetMatcherName(matcher Matcher, s string) {
	matcher.setMatcherName(s)
}

func targetMatch(matcher Matcher) bool {
	if !matcher.MatchDefaultTarget() {
		return true
	}
	if matcher.MatchOptionTarget() {
		return true
	}
	return false
}
