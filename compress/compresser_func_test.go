package compress

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"strings"
	"testing"
)

type MockFile struct {
	name  string
	size  int64
	isdir bool
	mode  os.FileMode
}

func (mock *MockFile) Name() string {
	return mock.name
}
func (mock *MockFile) Size() int64 {
	return mock.size
}
func (mock *MockFile) Mode() os.FileMode {
	return mock.mode
}
func (mock *MockFile) ModTime() time.Time {
	return time.Now()
}
func (mock *MockFile) IsDir() bool {
	return mock.isdir
}
func (mock *MockFile) Sys() interface{} {
	return nil
}

func CheckedMakeFile(file *os.File, create_file_name string) (flag bool) {
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

	if !CheckedMakeFile(fileio.file, create_file_name) {
		t.Errorf("got diff file name %s.", fileio.file.Name())
	}
	fileio.AllCloser()
	hostname = hostname + "/etc.tar.gz"
	os.Remove(hostname)
	hostname, _ = os.Hostname()
	hostname = "/mnt/" + hostname
	os.Remove(hostname)
}

func tmpWrite() {
	var (
		fileio           *Fileio = &Fileio{}		
		//*MockFileの定義はcompression_func_test.goにある
		mockfile *MockFile = &MockFile{name: "test.txt", size: 0, isdir: false, mode: os.ModePerm}
		mockgfile *MockFile = &MockFile{name: "gtest.txt", size: 9894688000, isdir: false, mode:os.ModePerm}
		mocks = []os.FileInfo {mockfile,mockgfile}
	)
	fileio.MakeFile("comp_test")
	fileio.CompressionFile(mocks, "comp_test")
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
		if hdr.Name != "comp_test/test.txt" && hdr.Name != "comp_test/gtest.txt"{
			t.Errorf("want comp_test/test.txt. got :%s\n", hdr.Name)
		}
	}
	os.Remove(remove_filename)
	os.Remove(hostname)
}
