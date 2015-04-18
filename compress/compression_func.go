package compress

import (
	"archive/tar"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var ed chan string = make(chan string)
var fin chan string = make(chan string)

func Compression(beforecheck_fileinfo []os.FileInfo, dirpath string, comfile *Fileio) {

	var (
		checked_fileinfo []os.FileInfo
		err              error
	)

	ChangeDir(dirpath)

L:
	for _, info := range beforecheck_fileinfo {		
		SetMatcherName(comfile, info.Name())
		if targetMatch(comfile) {
			continue L
		}
		ed <- info.Name()
		ef <- ""
		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			comfile.MakeFile(info.Name())
			evalsym, _ := os.Readlink(info.Name())
			hdr, _ := tar.FileInfoHeader(info, evalsym)
			hdr.Typeflag = tar.TypeSymlink
			if err = comfile.tw.WriteHeader(hdr); err != nil {
				fmt.Printf("write faild header symlink %s\n", info.Name())
				log.Fatal(err)
			}
			comfile.AllCloser()
		}
		if info.IsDir() {
			if checked_fileinfo, err = ioutil.ReadDir(info.Name()); err != nil {
				log.Fatal(err)
			}
			comfile.MakeFile(info.Name())
			if len(checked_fileinfo) != 0 {
				ChangeDir(info.Name())
				comfile.CompressionFile(checked_fileinfo, info.Name())
				comfile.AllCloser()
				ChangeDir(dirpath)
			} else {
				hdr, _ := tar.FileInfoHeader(info, "")
				hdr.Typeflag = tar.TypeDir
				if err = comfile.tw.WriteHeader(hdr); err != nil {
					fmt.Printf("write faild header symlink %s\n", info.Name())
					log.Fatal(err)
				}
				comfile.AllCloser()
			}
		}
	}
	fin <- "fin"
}
