package compress

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
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

func TestCompression(t *testing.T) {
	var (
		mockfile                                   *MockFile = &MockFile{name: "test.txt", size: 0, isdir: false, mode: os.ModePerm}
		mockdir                                    *MockFile = &MockFile{name: "test", size: 4096, isdir: true, mode: os.ModeDir}
		mockdir2                                   *MockFile = &MockFile{name: "test3", size: 4096, isdir: true, mode: os.ModeDir}
		mocklink                                   *MockFile = &MockFile{name: "test2", size: 17, isdir: false, mode: os.ModeSymlink}
		mocks                                                = []os.FileInfo{mockfile, mocklink, mockdir, mockdir2}
		dirpath, gopath, hostname, remove_filename string
		err                                        error
		check_file                                 *os.File
		hdr                                        *tar.Header
		fileReader                                 io.ReadCloser
	)
	hostname, _ = os.Hostname()
	hostname = "/mnt/" + hostname
	if _, err = os.Stat(hostname); err != nil {
		if err = os.Mkdir(hostname, os.ModePerm); err != nil {
			t.Errorf("can't start test")
		}
	}
	//GOPATHをとるための悪手
	o, _ := exec.Command(os.Getenv("SHELL"), "-c", "echo $GOPATH").Output()
	gopath = string(o)
	//stringsのTrimRightでchompのような動作
	gopath = strings.TrimRight(gopath, "\n")
	dirpath = gopath + "/src/suiage/compress/test"
	ChangeDir(dirpath)
	Compression(mocks, dirpath)

	hostname, _ = os.Hostname()
	hostname = "/mnt/" + hostname
	ChangeDir(hostname)
	for _, r := range []string{"/test.tar.gz", "/test2.tar.gz", "/test3.tar.gz"} {
		var buf bytes.Buffer
		remove_filename = hostname + r
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
			if hdr.Name != "test/hoge.txt" && hdr.Name != "test2" && hdr.Name != "test3/" {
				t.Errorf("Faild got :%s\n", hdr.Name)
			}
		}
		os.Remove(remove_filename)
	}
	os.Remove(hostname)

}
