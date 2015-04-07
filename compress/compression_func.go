package compress

import (
	"archive/tar"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func Compression(beforecheck_fileinfo []os.FileInfo, dirpath string) {

	var (
		checked_fileinfo []os.FileInfo
		err              error
		comfile          *Fileio = &Fileio{Target: &Target{}}
	)

	ChangeDir(dirpath)

L:
	for _, info := range beforecheck_fileinfo {
		SetMatcherName(comfile, info.Name())
		if targetMatch(comfile) {
			continue L
		}
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
}
