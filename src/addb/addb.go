package addb

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// TODO: Feature - ADDB indexes all entries. This lets entity ID to be independent of its path in ADDB.

type ADDB struct {
	Location string
	index    map[string]*ADDBComponent
}

func (db *ADDB) Init(addb_path string) error {
	if strings.HasPrefix(addb_path, "~") { // If path is relative to home directory
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		addb_path = home + strings.TrimPrefix(addb_path, "~")
	}
	addb_path = strings.TrimSuffix(addb_path, "/") // Remove trailing slash

	if _, err := os.Stat(addb_path); os.IsNotExist(err) {
		return errors.New("invalid ADDB path or directory not present")
	}
	db.Location = addb_path

	err := db.buildindex()
	if err != nil {
		return err
	}

	return nil
}

func (db *ADDB) GetComponent(id string) (*ADDBComponent, error) {

	entry := db.index[strings.TrimPrefix(id, "addb:")]
	if entry == nil {
		return nil, errors.New("cannot find '" + id + "' in ADDB.")
	}

	return entry, nil
	/*
		id = strings.TrimPrefix(id, "addb.") // remove "addb." from the path, if present
		id = strings.Replace(id, ".", "/", -1) // replace dots with slashes
		tempPath := db.Location + "/" + id + ".smspec"
		if _, err := os.Stat(tempPath); os.IsNotExist(err) {
			return nil, err
		}
		content, err := os.ReadFile(tempPath)
		if err != nil {
			return nil, err
		}

		var addb_component ADDBComponent
		err = yaml.Unmarshal(content, &addb_component)
		if err != nil {
			return nil, err
		}
		parts := strings.Split(id, "/")
		addb_component.Name = parts[len(parts) - 1]

		// update adm path
		var newPaths []string
		for _, adm_path := range addb_component.ADM {
			newPath := db.Location + "/" + id + "/" + adm_path
			newPaths = append(newPaths, newPath)
		}
		addb_component.ADM = newPaths

		switch (addb_component.Type) {
		case "human", "program", "system", "flow":
			return &addb_component, nil
		default:
			return nil, errors.New("unknown component type - " + string(addb_component.Type))
		}*/
}

////////////////////////////////////////
// Internal functions

func (db *ADDB) buildindex() error {
	if db.index == nil {
		db.index = make(map[string]*ADDBComponent)
	}

	files, err := traverse(db.Location, nil)
	if err != nil {
		return err
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		//var addb_component ADDBComponent
		//err = yaml.Unmarshal(content, &addb_component)
		components, err := unmarshalYamlBlocks(content)
		if err != nil {
			return err
		}

		for _, addb_component := range components {
			switch addb_component.Type {
			case "human", "program", "system", "flow":
				if _, present := db.index[addb_component.Id]; present {
					return errors.New("Found multiple entries in ADDB for '" + addb_component.Id + "'")
				}

				// Replace relative paths with absolute paths for ADM files
				var newPaths []string
				for _, adm_path := range addb_component.ADM {
					newPath := getBasePath(file) + "/" + adm_path
					newPaths = append(newPaths, newPath)
				}
				addb_component.ADM = newPaths

				db.index[addb_component.Id] = addb_component

			default:
				return errors.New("unknown component type - " + string(addb_component.Type))
			}
		}
	}

	return nil
}

func unmarshalYamlBlocks(content []byte) ([]*ADDBComponent, error) {
	var out []*ADDBComponent
	read := bytes.NewReader(content)
	decoder := yaml.NewDecoder(read)
	for {
		var addb_component ADDBComponent
		if err := decoder.Decode(&addb_component); err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		out = append(out, &addb_component)
	}
	return out, nil
}

func getBasePath(filePath string) string {
	parts := strings.Split(filePath, "/")
	basePath := strings.Join(parts[:len(parts)-1], "/")
	return basePath
}

////////////////////////////////////////
// Helper functions

func traverse(path string, ignoreList []string) ([]string, error) {
	var files, ignore []string

	ignore = append(ignore, ignoreList...)

	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if fileInfo.IsDir() {
		items, _ := ioutil.ReadDir(path)
		// first process .git and .gitignore
		for _, item := range items {
			if item.Name() == ".git" { // don't process '.git' folder
				ignore = append(ignore, ".git")
				continue
			} else if item.Name() == ".gitignore" { // skip everything listed in .gitignore
				ignore = append(ignore, ".gitignore")
				content, err := os.ReadFile(path + "/.gitignore")
				if err != nil {
					return nil, err
				}
				ignoreFiles := string(content)
				ignore = append(ignore, strings.Split(ignoreFiles, "\n")...)
			} else { // Only pick files with extension .smspec
				fileparts := strings.Split(item.Name(), ".")
				fileExtension := fileparts[len(fileparts)-1]
				if contains(item.Name(), ignore) {
					continue
				}
				if item.IsDir() { // subdirectories
					f, err := traverse(path+"/"+item.Name(), ignore)
					if err != nil {
						return nil, err
					}
					files = append(files, f...)
				} else if fileExtension != "smspec" {
					continue
				} else {
					files = append(files, path+"/"+item.Name())
				}
			}
		}
	} else { // received a single file's path
		files = append(files, path)
	}

	return files, nil
}

func contains(item string, list []string) bool {
	for _, x := range list {
		if x == item {
			return true
		}
	}

	return false
}
