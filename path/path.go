package path

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// Dirs ...
func Dirs(dirname string, dirs *[]string) (err error) {
	var fileInfos []os.FileInfo
	fileInfos, err = ioutil.ReadDir(dirname)
	if err != nil {
		return
	}
	for _, fi := range fileInfos {
		if fi.IsDir() {
			dir := filepath.Join(dirname, fi.Name())
			*dirs = append(*dirs, dir)
			err = Dirs(dir, dirs)
			if err != nil {
				return
			}
		}
	}
	return
}
