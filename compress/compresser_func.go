package compress

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Compresser interface {
	CompressionFile([]os.FileInfo, string)
	MakeFile(string)
	AllCloser()
}

type Fileio struct {
	fileWriter io.WriteCloser
	tw         *tar.Writer
	file       *os.File
	*Target
}

func (comfile *Fileio) AllCloser() {
	defer comfile.file.Close()
	defer comfile.fileWriter.Close()
	defer comfile.tw.Close()
}

func (comfile *Fileio) MakeFile(create_file_name string) {
	var (
		hostname                     string
		err                          error
		year, day                    int
		month                        time.Month
		str_year, str_month, str_day string
	)
	if hostname, err = os.Hostname(); err != nil {
		log.Fatal(err)
	}
	if create_file_name == "" {
		year, month, day = time.Now().Date()
		str_year = strconv.Itoa(year)
		str_month = strconv.Itoa(int(month))
		str_day = strconv.Itoa(day)
		hostname = "/mnt/" + hostname + "_" + str_year + "_" + str_month + "_" + str_day + ".tar.gz"
	} else {
		hostname = "/mnt/" + hostname + "/" + create_file_name + ".tar.gz"
	}
	if comfile.file, err = os.Create(hostname); err != nil {
		log.Fatal(err)
	}
	fmt.Println(hostname)
	comfile.fileWriter = comfile.file
	comfile.fileWriter = gzip.NewWriter(comfile.file)
	comfile.tw = tar.NewWriter(comfile.fileWriter)
}

func (f *Fileio) CompressionFile(checked_fileinfo []os.FileInfo, dirname string) {
	var (
		err            error
		tmp_fileinfo   []os.FileInfo
		change_dirpath string
	)
compress:
	for _, infile := range checked_fileinfo {
		var target *Target = &Target{filepath.Join(dirname, infile.Name())}
		if infile.IsDir() {
			if target.MatchOptionTarget() {
				continue compress
			}
			if tmp_fileinfo, err = ioutil.ReadDir(infile.Name()); err != nil {
				log.Fatal(err)
			}
			tmpname := filepath.Join(dirname, infile.Name())
			hdr, _ := tar.FileInfoHeader(infile, "")
			hdr.Typeflag = tar.TypeDir
			hdr.Name = tmpname
			if err = f.tw.WriteHeader(hdr); err != nil {
				fmt.Printf("write faild header Dir %s\n", tmpname)
				log.Fatal(err)
			}
			change_dirpath, _ = filepath.Abs(infile.Name())
			ChangeDir(change_dirpath)
			dirname = filepath.Join(dirname, infile.Name())
			f.CompressionFile(tmp_fileinfo, dirname)
			dirname, _ = filepath.Split(dirname)
			change_dirpath, _ = filepath.Split(change_dirpath)
			ChangeDir(change_dirpath)
			tmp_fileinfo = nil
		} else {
			tmpname := filepath.Join(dirname, infile.Name())
			target.name = tmpname
			if target.MatchOptionTarget() {
				continue compress
			}
			if infile.Mode()&os.ModeSymlink == os.ModeSymlink {
				evalsym, _ := os.Readlink(infile.Name())
				hdr, _ := tar.FileInfoHeader(infile, evalsym)
				hdr.Typeflag = tar.TypeSymlink
				hdr.Name = tmpname
				if err = f.tw.WriteHeader(hdr); err != nil {
					fmt.Printf("write faild header symlink %s\n", tmpname)
					log.Fatal(err)
				}
			} else {
				body, _ := ioutil.ReadFile(infile.Name())
				hdr, _ := tar.FileInfoHeader(infile, "")
				hdr.Typeflag = tar.TypeReg
				hdr.Name = tmpname
				if err = f.tw.WriteHeader(hdr); err != nil {
					fmt.Printf("write faild header %s\n", tmpname)
					log.Fatal(err)
				}
				if body != nil {
					f.tw.Write(body)
				}
			}
		}
	}
}
