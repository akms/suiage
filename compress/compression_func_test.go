package compress

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
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
	default_except_targets = strings.Fields(`^lost\+found$ ^proc$ ^sys$ ^dev$ ^mnt$ ^media$ ^run$ ^selinux$ ^run$`)
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

func TestCompressionFile(t *testing.T) {
	var (
		br               *bytes.Reader
		gw               *gzip.Writer
		gr               *gzip.Reader
		tw               *tar.Writer
		tr               *tar.Reader
		file, check_file *os.File
		file_stat                 os.FileInfo
		hostname, remove_filename string
		checked_fileinfo          []os.FileInfo
		err                       error
		hdr                       *tar.Header
	)

	hostname, _ = os.Hostname()
	hostname = "/mnt/" + hostname
	os.Mkdir(hostname, os.ModePerm)

	option_except_targets = strings.Fields(`^lost\+found$ ^proc$ ^sys$ ^dev$ ^mnt$ ^media$ ^run$ ^selinux$ ^run$ ^tmp$ ^_old$ ^boot$ ^opt$ ^root$ ^sbin$ ^etc$ ^var$ ^home$`)

	ChangeDir("/")
	checked_fileinfo, _ = ioutil.ReadDir("/srv")
	gw, tw, file = MakeFile("srv")
	CompressionFile(tw, checked_fileinfo, "srv")
	defer file.Close()
	defer gw.Close()
	defer tw.Close()

	remove_filename = hostname + "/srv.tar.gz"
	body, _ := ioutil.ReadFile(remove_filename)
	check_file,_ = os.Open(remove_filename)
	br = bytes.NewReader(body)
	
	file_stat,_ = check_file.Stat()
	hdr,_ = tar.FileInfoHeader(file_stat,"")
	fmt.Println(hdr)
	fmt.Println(br)
	gr, err = gzip.NewReader(check_file)
	if err != nil {
		t.Errorf("Can't read gr")
	}
	
	tr = tar.NewReader(check_file)
	
	for {
		hdr, err = tr.Next()

		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("Can't read hdr %s\n", err)
			break
		} else {
			fmt.Printf("%s:\n", hdr.Name)
		}
	}
	os.Remove(remove_filename)
	os.Remove(hostname)
	defer check_file.Close()
	defer gr.Close()
}
