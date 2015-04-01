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
	return
}

func TestMakeFile(t *testing.T) {
	var (
		fileio           *Fileio = &Fileio{}
		create_file_name string
		hostname         string
	)

	fileio.MakeFile("")

	if fileio.fileWriter == nil {
		t.Errorf("make faild 1st gzip writer.")
	}
	if fileio.tw == nil {
		t.Errorf("make faild 1st tar writer.")
	}
	if fileio.file == nil {
		t.Errorf("make faild 1st file.")
	}

	if !NoNameMakeFile(fileio.file) {
		t.Errorf("got diff file name %s.", fileio.file.Name())
	}
	fileio.AllCloser()

	hostname, _ = os.Hostname()
	hostname = "/mnt/" + hostname
	os.Mkdir(hostname, os.ModePerm)
	create_file_name = "etc"

	fileio.MakeFile(create_file_name)

	if fileio.fileWriter == nil {
		t.Errorf("make faild 2nd gzip writer.")
	}
	if fileio.tw == nil {
		t.Errorf("make faild 2nd tar writer.")
	}
	if fileio.file == nil {
		t.Errorf("make faild 2nd file.")
	}

	if !NamedMakeFile(fileio.file, create_file_name) {
		t.Errorf("got diff file name %s.", fileio.file.Name())
	}
	fileio.AllCloser()

	var c Compresser = &Fileio{}
	c.MakeFile(create_file_name)
	c.AllCloser()
	hostname, _ = os.Hostname()
	hostname = "/mnt/" + hostname
	os.Remove(hostname)
}

func tmpWrite() {
	var (
		file             *os.File
		fileio           *Fileio = &Fileio{}
		checked_fileinfo []os.FileInfo
	)
	ChangeDir("/tmp")
	os.Mkdir("comp_test", os.ModePerm)
	ChangeDir("comp_test")
	file, _ = os.Create("test.txt")
	defer file.Close()

	ChangeDir("/tmp")
	checked_fileinfo, _ = ioutil.ReadDir("comp_test")
	fileio.MakeFile("comp_test")
	ChangeDir("/tmp/comp_test")
	fileio.CompressionFile(checked_fileinfo, "comp_test")
	fileio.AllCloser()
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
	option_except_targets = strings.Fields(`^lost\+found$ ^proc$ ^sys$ ^dev$ ^mnt$ ^media$ ^run$ ^selinux$ ^tmp$ ^_old$ ^boot$ ^opt$ ^root$ ^sbin$ ^etc$ ^var$ ^home$ ^srv$`)

	tmpWrite()

	//以下今回のテストの目的である.tar.gzファイルの読み込み

	remove_filename = hostname + "/comp_test.tar.gz"
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
		if hdr.Name != "comp_test/test.txt" {
			t.Errorf("want comp_test/test.txt. got :%s\n", hdr.Name)
		}
	}
	os.Remove(remove_filename)
	os.Remove("/tmp/comp_test/test.txt")
	os.Remove("/tmp/comp_test")
	os.Remove(hostname)
}
