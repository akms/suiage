package suiage

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	gw                              *gzip.Writer
	tw                              *tar.Writer
	file                            *os.File
	default_targets, option_targets []string
)

func MakeFile() (*gzip.Writer, *tar.Writer, *os.File) {
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
	year, month, day = time.Now().Date()
	str_year = strconv.Itoa(year)
	str_month = strconv.Itoa(int(month))
	str_day = strconv.Itoa(day)
	hostname = "/mnt/" + hostname + "_" + str_year + "_" + str_month + "_" + str_day + ".tar.gz"
	if file, err = os.Create(hostname); err != nil {
		log.Fatal(err)
	}
	gw = gzip.NewWriter(file)
	tw = tar.NewWriter(gw)
	return gw, tw, file
}

func MatchDefaultTarget(name string) bool {
	for i, s := range default_targets {
		default_Regexp := regexp.MustCompile(s)
		if default_Regexp.MatchString(name) {
			default_targets = append(default_targets[:i], default_targets[i+1:]...)
			return false
		}
	}
	return true
}

func MatchOptionTarget(name string) bool {
	for _, s := range option_targets {
		option_Regexp := regexp.MustCompile(s)
		if option_Regexp.MatchString(name) {
			//option_targets = append(option_targets[:i], option_targets[i+1:]...)
			fmt.Println(name)
			fmt.Println(option_targets)
			return true
		}
	}
	return false
}

func CheckTarget(dirpath string) {
	var (
		beforecheck_fileinfo, checked_fileinfo []os.FileInfo
		err                                    error
	)
	default_targets = strings.Fields(`^lost\+found$ ^proc$ ^sys$ ^dev$ ^mnt$ ^media$ ^run$ ^selinux$`)
	option_targets = ReadOption()
	ChangeDir(dirpath)
	if beforecheck_fileinfo, err = ioutil.ReadDir(dirpath); err != nil {
		log.Fatal(err)
	}
	for _, info := range beforecheck_fileinfo {
		if MatchDefaultTarget(info.Name()) {
			checked_fileinfo = append(checked_fileinfo, info)
		}
	}
	_, dirname := filepath.Split(dirpath)
	gw, tw, file = MakeFile()
	CompressionFile(tw, checked_fileinfo, dirname)
	defer file.Close()
	defer gw.Close()
	defer tw.Close()
}

/*func walkFn(path string,info os.FileInfo,err error) error {
	fmt.Println(filepath.Base(info.Name()))
	return nil
}
if err = infile.Walk(file.Name(),walkFn); err != nil {
	log.Fatal(err)
}*/

func CompressionFile(tw *tar.Writer, checked_fileinfo []os.FileInfo, dirname string) {
	var (
		err            error
		tmp_fileinfo   []os.FileInfo
		change_dirpath string
	)
compress:
	for _, infile := range checked_fileinfo {
		if infile.IsDir() {
			target_name := filepath.Join(dirname, infile.Name())
			if MatchOptionTarget(target_name) {
				continue compress
			}
			if tmp_fileinfo, err = ioutil.ReadDir(infile.Name()); err != nil {
				log.Fatal(err)
			}
			change_dirpath, _ = filepath.Abs(infile.Name())
			ChangeDir(change_dirpath)
			dirname = filepath.Join(dirname, infile.Name())
			CompressionFile(tw, tmp_fileinfo, dirname)
			dirname, _ = filepath.Split(dirname)
			change_dirpath, _ = filepath.Split(change_dirpath)
			ChangeDir(change_dirpath)
		} else {
			tmpname := filepath.Join(dirname, infile.Name())
			if MatchOptionTarget(tmpname) {
				continue compress
			}
			if infile.Mode()&os.ModeSymlink == os.ModeSymlink {
				evalsym, _ := os.Readlink(infile.Name())
				hdr, _ := tar.FileInfoHeader(infile, evalsym)
				hdr.Typeflag = tar.TypeSymlink
				hdr.Name = tmpname
				if err = tw.WriteHeader(hdr); err != nil {
					log.Fatal(err)
				}
			} else {
				//fmt.Println(tmpname)
				body, _ := ioutil.ReadFile(infile.Name())
				hdr, _ := tar.FileInfoHeader(infile, "")
				hdr.Typeflag = tar.TypeRegA
				hdr.Name = tmpname
				if err = tw.WriteHeader(hdr); err != nil {
					log.Fatal(err)
				}
				if _, err = tw.Write(body); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
