package compress

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
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
		//	gw               *gzip.Writer
		fileWriter       io.WriteCloser
		tw               *tar.Writer
		file             *os.File
		create_file_name string
		hostname         string
	)
	//gw, tw, file = MakeFile("")
	fileWriter, tw, file = MakeFile("")
	//if gw == nil {
	if fileWriter == nil {
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
	//defer gw.Close()
	defer fileWriter.Close()
	defer tw.Close()

	hostname, _ = os.Hostname()
	hostname = "/mnt/" + hostname
	os.Mkdir(hostname, os.ModePerm)
	create_file_name = "etc"
	//	gw, tw, file = MakeFile(create_file_name)
	fileWriter, tw, file = MakeFile(create_file_name)

	//	if gw == nil {
	if fileWriter == nil {
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
	//defer gw.Close()
	defer fileWriter.Close()
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
	//str := strings.Fields("_old boot opt root sbin etc var home")
	str := strings.Fields("_old")
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

func tmpWrite() {
	var (
		fileWriter       io.WriteCloser
		tw               *tar.Writer
		file             *os.File
		checked_fileinfo []os.FileInfo
	)
	ChangeDir("/")
	checked_fileinfo, _ = ioutil.ReadDir("/srv")
	fileWriter, tw, file = MakeFile("srv")
	ChangeDir("/srv")
	CompressionFile(tw, checked_fileinfo, "srv")
	defer file.Close()
	defer fileWriter.Close()
	defer tw.Close()
}

func TestCompressionFile(t *testing.T) {
	var (
		check_file                *os.File
		hostname, remove_filename string
		err                       error
		hdr                       *tar.Header
		fileReader                io.ReadCloser
		buf                       bytes.Buffer
	)

	hostname, _ = os.Hostname()
	hostname = "/mnt/" + hostname
	os.Mkdir(hostname, os.ModePerm)
	option_except_targets = strings.Fields(`^lost\+found$ ^proc$ ^sys$ ^dev$ ^mnt$ ^media$ ^run$ ^selinux$ ^tmp$ ^_old$ ^boot$ ^opt$ ^root$ ^sbin$ ^etc$ ^var$ ^home$`)

	tmpWrite()

	//以下今回のテストの目的である.tar.gzファイルの読み込み
	//.tarファイルでは動作することが確認できている
	remove_filename = hostname + "/srv.tar.gz"
	ChangeDir(hostname)
	check_file, err = os.Open(remove_filename)
	if err != nil {
		t.Errorf("Can't open file %s\n", remove_filename)
	}
	defer check_file.Close()

	_, err = io.Copy(&buf, check_file)
	if fileReader, err = gzip.NewReader(&buf); err != nil {
		t.Errorf("%s", err)
	}
	defer fileReader.Close()
	tr := tar.NewReader(fileReader)
	for {
		hdr, err = tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("Can't read hdr %s\n", err)
			break
		}
		if hdr.Name != "srv/test.txt" {
			t.Errorf("want srv/test.txt. got :%s\n",hdr.Name)
		}
	}
	os.Remove(remove_filename)
	os.Remove(hostname)
}
