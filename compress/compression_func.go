package compress

import (
	"archive/tar"
	//	"bytes"
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

var (
	fileWriter                                    io.WriteCloser
	tw                                            *tar.Writer
	file                                          *os.File
	default_except_targets, option_except_targets []string
)

func MakeFile(create_file_name string) (io.WriteCloser, *tar.Writer, *os.File) {
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
	if file, err = os.Create(hostname); err != nil {
		log.Fatal(err)
	}
	fmt.Println(hostname)
	fileWriter = file
	fileWriter = gzip.NewWriter(file)
	tw = tar.NewWriter(fileWriter)
	return fileWriter, tw, file
}

func MatchDefaultTarget(name string) bool {
	for i, s := range default_except_targets {
		default_Regexp := regexp.MustCompile(s)
		if default_Regexp.MatchString(name) {
			default_except_targets = append(default_except_targets[:i], default_except_targets[i+1:]...)
			return false
		}
	}
	return true
}

func MatchOptionTarget(name string) bool {
	for _, s := range option_except_targets {
		option_Regexp := regexp.MustCompile(s)
		if option_Regexp.MatchString(name) {
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
	default_except_targets = strings.Fields(`^lost\+found$ ^proc$ ^sys$ ^dev$ ^mnt$ ^media$ ^run$ ^selinux$ ^_old$`)
	option_except_targets = ReadOption()
	ChangeDir(dirpath)
	if beforecheck_fileinfo, err = ioutil.ReadDir(dirpath); err != nil {
		log.Fatal(err)
	}
L:
	for _, info := range beforecheck_fileinfo {
		if MatchDefaultTarget(info.Name()) {
			if MatchOptionTarget(info.Name()) {
				continue L
			}
			if info.Mode()&os.ModeSymlink == os.ModeSymlink {
				tmpname := filepath.Join(dirpath, info.Name())
				fileWriter, tw, file = MakeFile(info.Name())
				evalsym, _ := os.Readlink(info.Name())
				hdr, _ := tar.FileInfoHeader(info, evalsym)
				hdr.Typeflag = tar.TypeSymlink
				//hdr.Name = tmpname
				if err = tw.WriteHeader(hdr); err != nil {
					fmt.Printf("write faild header symlink %s\n", tmpname)
					log.Fatal(err)
				}
				defer file.Close()
				defer fileWriter.Close()
				defer tw.Close()
			}
			if info.IsDir() {
				if checked_fileinfo, err = ioutil.ReadDir(info.Name()); err != nil {
					log.Fatal(err)
				}
				ChangeDir(info.Name())
				fileWriter, tw, file = MakeFile(info.Name())
				CompressionFile(tw, checked_fileinfo, info.Name())
				defer file.Close()
				defer fileWriter.Close()
				defer tw.Close()
				ChangeDir(dirpath)
			}
		}
	}
}

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
			if len(tmp_fileinfo) == 0 {
				tmpname := filepath.Join(dirname, infile.Name())
				hdr, _ := tar.FileInfoHeader(infile, "")
				hdr.Typeflag = tar.TypeDir
				hdr.Name = tmpname
				if err = tw.WriteHeader(hdr); err != nil {
					fmt.Printf("write faild header Dir %s\n", tmpname)
					log.Fatal(err)
				}
				continue compress
			} else {
				change_dirpath, _ = filepath.Abs(infile.Name())
				ChangeDir(change_dirpath)
				dirname = filepath.Join(dirname, infile.Name())
				CompressionFile(tw, tmp_fileinfo, dirname)
				dirname, _ = filepath.Split(dirname)
				change_dirpath, _ = filepath.Split(change_dirpath)
				ChangeDir(change_dirpath)
			}
			tmp_fileinfo = nil
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
					fmt.Printf("write faild header symlink %s\n", tmpname)
					log.Fatal(err)
				}
			} else {
				//io.Copyを使用しない場合はos.Openからioutil.ReadFileに変更
				//body, _ := os.Open(infile.Name())
				body, _ := ioutil.ReadFile(infile.Name())
				hdr, _ := tar.FileInfoHeader(infile, "")
				//180Mのバイナリファイルでwrite too long エラーが出たので
				//tar.TypeRegAからtar.TypeRegへ変更
				//hdr.Typeflag = tar.TypeRegA
				hdr.Typeflag = tar.TypeReg
				hdr.Name = tmpname
				if err = tw.WriteHeader(hdr); err != nil {
					fmt.Printf("write faild header %s\n", tmpname)
					log.Fatal(err)
				}
				if body != nil {
					//io.Copyでは動作が安定せず、途中で書き込みに失敗する
					//要原因究明
					/*var buf bytes.Buffer

					if _, err = io.Copy(&buf, body); err != nil {
						fmt.Printf("write faild %s\n",tmpname)
						log.Fatal(err)
					}
					tw.Write(buf.Bytes())*/
					tw.Write(body)
				}
			}
		}
	}
}
