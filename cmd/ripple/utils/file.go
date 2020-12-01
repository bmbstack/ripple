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
