package files

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func YamlFromFilePathToModel(filePath string, model interface{}) (err error) {
	file, _ := ioutil.ReadFile(filePath)
	err = yaml.Unmarshal(file, model)
	return
}
