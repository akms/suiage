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
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Compresser interface {
	CompressionFile([]os.FileInfo, string)
	MakeFile(string)
	AllCloser()
}

type Matcher interface {
	MatchDefaultTarget() bool
	MatchOptionTarget() bool
}

type Target struct {
	name string
}

type Fileio struct {
	fileWriter io.WriteCloser
	tw         *tar.Writer
	file       *os.File
}

var (
	default_except_targets, option_except_targets []string = strings.Fields(`^lost\+found$ ^proc$ ^sys$ ^dev$ ^mnt$ ^media$ ^run$ ^selinux$ ^boot$ ^_old$`), ReadOption()
)

func (comfile *Fileio) AllCloser() {
	defer comfile.file.Close()
	defer comfile.fileWriter.Close()
	defer comfile.tw.Close()
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

func (target *Target)MatchOptionTarget() bool {
	for _, s := range option_except_targets {
		option_Regexp := regexp.MustCompile(s)
		if option_Regexp.MatchString(target.name) {
			return true
		}
	}
	return false
}

func CheckTarget(dirpath string) {
	var (
		beforecheck_fileinfo, checked_fileinfo []os.FileInfo
		err                                    error
		comfile                                *Fileio = &Fileio{}
	)
	ChangeDir(dirpath)
	if beforecheck_fileinfo, err = ioutil.ReadDir(dirpath); err != nil {
		log.Fatal(err)
	}
L:
	for _, info := range beforecheck_fileinfo {
		var target *Target = &Target{info.Name()}
		if target.MatchDefaultTarget() {
			if target.MatchOptionTarget() {
				continue L
			}
			if info.Mode()&os.ModeSymlink == os.ModeSymlink {
				tmpname := filepath.Join(dirpath, info.Name())
				comfile.MakeFile(info.Name())
				evalsym, _ := os.Readlink(info.Name())
				hdr, _ := tar.FileInfoHeader(info, evalsym)
				hdr.Typeflag = tar.TypeSymlink
				if err = comfile.tw.WriteHeader(hdr); err != nil {
					fmt.Printf("write faild header symlink %s\n", tmpname)
					log.Fatal(err)
				}
				comfile.AllCloser()
			}
			if info.IsDir() {
				if checked_fileinfo, err = ioutil.ReadDir(info.Name()); err != nil {
					log.Fatal(err)
				}
				if len(checked_fileinfo) != 0 {
					ChangeDir(info.Name())
					comfile.MakeFile(info.Name())
					comfile.CompressionFile(checked_fileinfo, info.Name())
					comfile.AllCloser()
					ChangeDir(dirpath)
				} else {
					tmpname := filepath.Join(dirpath, info.Name())
					comfile.MakeFile(info.Name())
					hdr, _ := tar.FileInfoHeader(info, "")
					hdr.Typeflag = tar.TypeDir
					if err = comfile.tw.WriteHeader(hdr); err != nil {
						fmt.Printf("write faild header symlink %s\n", tmpname)
						log.Fatal(err)
					}
					comfile.AllCloser()
				}
			}
		}
	}
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
