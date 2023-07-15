package args

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Common 'command' interface
// NOTE: Not used currently. All commands' execute function are directly called without using an interface.
//type command interface {
//	execute() error
//}

////////////////////////////////////////
// Helper structure to collect all files in directory and its sub-directories
type filelist struct {
	files []string
}

func (l *filelist) WalkDir(path string, d fs.DirEntry, err error) error {
	if (d.IsDir()) {
		return nil
	} else { // is a file
		if filepath.Ext(path) == ".smspec" {
			l.files = append(l.files, path)
		}
	}
	return nil
}

////////////////////////////////////////
// Common functions used across the package

func getContent(path string) (map[string]string, error) {
	fileAndContent := make(map[string]string)

	files, err := getFiles(path)
	if err != nil {
		return nil, err
	}
	if len(files) > 0 {
		fmt.Println("Found", len(files), "file(s)")
	}
	for _, file := range files {
		newContent, err := getFileContent(file)
		if err != nil {
			return nil, err
		}
		fileAndContent[getFileName(file)] = newContent
	}

	return fileAndContent, nil
}

func getFiles(path string) ([]string, error) {
	/*fileInfo*/_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	var list filelist
	filepath.WalkDir(path, list.WalkDir)
	return list.files, nil
}

func getFileContent(path string) (string, error) {
	var content string
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		content = content + "\n" + line
	}

	return content, nil
}

func getFileName(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts) - 1]
}

func checkAndCreateDirectory(directory string) string {
	if directory[len(directory) - 1] != '/' { // append a "/" if directory string doesn't contain it.
		directory = directory + "/"
	}
	var path string
	for _, dir := range strings.Split(directory, "/") {
		if dir != "" {
			path += dir
			if _, err := os.Stat(path); err != nil {
				err = os.Mkdir(path, 0700) // create directory with RW rights to owner only.
				if err != nil {
					panic(err)
				}
			}
			path += "/"
		}
	}

	return directory
}

func PrintErrors(errs []error) {
	for _, err := range errs {
		fmt.Println("ERROR: " + err.Error())
	}
}