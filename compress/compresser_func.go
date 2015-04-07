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

func (f *Fileio) AllCloser() {
	defer f.file.Close()
	defer f.fileWriter.Close()
	defer f.tw.Close()
}

func (f *Fileio) MakeFile(create_file_name string) {
	var (
		hostname string
		err      error
	)
	if hostname, err = os.Hostname(); err != nil {
		log.Fatal(err)
	}
	if create_file_name == "" {
		err = fmt.Errorf("create_file_name is nil")
		log.Fatal(err)
	} else {
		hostname = "/mnt/" + hostname + "/" + create_file_name + ".tar.gz"
	}
	if f.file, err = os.Create(hostname); err != nil {
		log.Fatal(err)
	}
	fmt.Println(hostname)
	f.fileWriter = gzip.NewWriter(f.file)
	f.tw = tar.NewWriter(f.fileWriter)
}

func (f *Fileio) CompressionFile(checked_fileinfo []os.FileInfo, dirname string) {
	var (
		err            error
		tmp_fileinfo   []os.FileInfo
		change_dirpath string
	)
	f.Target = &Target{}
compress:
	for _, infile := range checked_fileinfo {
		tmpname := filepath.Join(dirname, infile.Name())
		SetMatcherName(f, tmpname)
		if targetMatch(f) {
			continue compress
		}
		fmt.Println(tmpname)
		if infile.IsDir() {
			if tmp_fileinfo, err = ioutil.ReadDir(infile.Name()); err != nil {
				log.Fatal(err)
			}
			hdr, _ := tar.FileInfoHeader(infile, "")
			hdr.Typeflag = tar.TypeDir
			hdr.Name = tmpname
			fmt.Println(tmpname)
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
			if infile.Mode()&os.ModeSymlink == os.ModeSymlink {
				evalsym, _ := os.Readlink(infile.Name())
				hdr, _ := tar.FileInfoHeader(infile, evalsym)
				hdr.Typeflag = tar.TypeSymlink
				hdr.Name = tmpname
				fmt.Println(tmpname)
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
