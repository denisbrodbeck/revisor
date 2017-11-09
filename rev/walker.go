package rev

import (
	"log"
	"os"
	"path/filepath"
)

// Find takes a root path and returns all files underneath root path.
// This includes files in sub-directories.
func Find(root, baseURL string) (Files, error) {
	files := []File{}
	if err := filepath.Walk(root, walker(root, baseURL, &files)); err != nil {
		return nil, err
	}
	return files, nil
}

func walker(root, baseURL string, files *[]File) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				log.Println("Insufficient permissions:", err)
			} else {
				return err
			}
		}
		// skip symlinks, sockets, etc.
		if info.Mode().IsRegular() == false {
			return nil
		}
		// skip directories
		if info.IsDir() {
			return nil
		}
		file, err := NewFile(path, root, baseURL, info.Size(), info.Mode())
		if err != nil {
			return err
		}

		*files = append(*files, *file)

		return nil
	}
}
