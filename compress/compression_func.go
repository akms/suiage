package compress

import (
	"archive/tar"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

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
