package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func ReadConfig(cfg interface{}, path string, fileName string) error {

	getFormatFile := filePath(path)

	switch getFormatFile {
	case ".json":
		fname := fmt.Sprintf("%s/%s", path, fileName)
		jsonFile, err := ioutil.ReadFile(fname)
		if err != nil {
			return err
		}
		return json.Unmarshal(jsonFile, cfg)
	default:
		fname := fmt.Sprintf("%s/%s", path, fileName)
		yamlFile, err := ioutil.ReadFile(fname)
		if err != nil {
			return err
		}
		return yaml.Unmarshal(yamlFile, cfg)
	}

}

func filePath(root string) string {
	var file string
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		file = filepath.Ext(info.Name())
		return nil
	})
	return file
}
