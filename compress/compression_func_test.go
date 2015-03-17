package compress

import (
	"archive/tar"
	"compress/gzip"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func NoNameMakeFile(file *os.File) (flag bool) {
	var (
		remove_filename              string
		hostname                     string
		year, day                    int
		month                        time.Month
		str_year, str_month, str_day string
	)
	hostname, _ = os.Hostname()
	year, month, day = time.Now().Date()
	str_year = strconv.Itoa(year)
	str_month = strconv.Itoa(int(month))
	str_day = strconv.Itoa(day)
	hostname = "/mnt/" + hostname + "_" + str_year + "_" + str_month + "_" + str_day + ".tar.gz"
	remove_filename = file.Name()
	if hostname != remove_filename {
		flag = false
	} else {
		flag = true
	}
	os.Remove(remove_filename)
	return
}

func NamedMakeFile(file *os.File, create_file_name string) (flag bool) {
	var (
		remove_filename string
		hostname        string
	)
	hostname, _ = os.Hostname()
	hostname = "/mnt/" + hostname + "/" + create_file_name + ".tar.gz"
	remove_filename = file.Name()
	if hostname != remove_filename {
		flag = false
	} else {
		flag = true
	}
	os.Remove(remove_filename)
	hostname, _ = os.Hostname()
	hostname = "/mnt/" + hostname
	os.Remove(hostname)
	return
}

func TestMakeFile(t *testing.T) {
	var (
		gw               *gzip.Writer
		tw               *tar.Writer
		file             *os.File
		create_file_name string
		hostname         string
	)
	gw, tw, file = MakeFile("")
	if gw == nil {
		t.Errorf("make faild 1st gzip writer.")
	}
	if tw == nil {
		t.Errorf("make faild 1st tar writer.")
	}
	if file == nil {
		t.Errorf("make faild 1st file.")
	}

	if !NoNameMakeFile(file) {
		t.Errorf("got file name %s.", file.Name())
	}
	defer file.Close()
	defer gw.Close()
	defer tw.Close()

	hostname, _ = os.Hostname()
	hostname = "/mnt/" + hostname
	os.Mkdir(hostname, os.ModePerm)
	create_file_name = "etc"
	gw, tw, file = MakeFile(create_file_name)

	if gw == nil {
		t.Errorf("make faild 2nd gzip writer.")
	}
	if tw == nil {
		t.Errorf("make faild 2nd tar writer.")
	}
	if file == nil {
		t.Errorf("make faild 2nd file.")
	}

	if !NamedMakeFile(file, create_file_name) {
		t.Errorf("got file name %s.", file.Name())
	}
	defer file.Close()
	defer gw.Close()
	defer tw.Close()
}

func TestMatchDefaultTarget(t *testing.T) {
	default_except_targets = strings.Fields(`^lost\+found$ ^proc$ ^sys$ ^dev$ ^mnt$ ^media$ ^run$ ^selinux$`)
	str := strings.Fields("lost+found proc sys dev mnt media selinux run")
	for _, s := range str {
		if MatchDefaultTarget(s) {
			t.Errorf("Match faild %s", s)
		}
	}
	default_except_targets = strings.Fields(`^lost\+found$ ^proc$ ^sys$ ^dev$ ^mnt$ ^media$ ^run$ ^selinux$ ^run`)
	str = strings.Fields("var etc lib lib64")
	for _, s := range str {
		if !MatchDefaultTarget(s) {
			t.Errorf("Match faild %s", s)
		}
	}
}

func TestMatchOptionTarget(t *testing.T) {
	str := strings.Fields("_old boot opt root sbin etc var home")
	option_except_targets = ReadOption()
	for _, s := range str {
		if !MatchOptionTarget(s) {
			t.Errorf("Match faild %s. Please check /etc/suiage.conf", s)
		}
	}
	str = strings.Fields("lost+found proc sys dev mnt media selinux run")
	for _, s := range str {
		if MatchOptionTarget(s) {
			t.Errorf("Match faild %s", s)
		}
	}
}
