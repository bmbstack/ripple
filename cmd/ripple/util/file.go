package utils

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

// CollectFiles Collect the files with these extensions under src
func CollectFiles(dir string, extensions []string) ([]string, error) {
	files := []string{}
	err := filepath.Walk(dir, func(file string, info os.FileInfo, err error) error {
		// If we have an err pass it up
		if err != nil {
			return err
		}
		// Deal with files only
		if !info.IsDir() {
			// Check for go files
			name := path.Base(file)
			if !strings.HasPrefix(name, ".") {
				for _, extension := range extensions {
					if strings.HasSuffix(name, extension) {
						files = append(files, file)
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return files, err
	}
	return files, nil
}

// Exist 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
