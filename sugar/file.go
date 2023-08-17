package sugar

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

// WriteFile 写文件
func WriteFile(file string, data []byte) error {
	path := filepath.Dir(file)

	if !IsDir(path) {
		err := os.MkdirAll(path, fs.ModePerm)
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(file, data, fs.ModePerm)
}
