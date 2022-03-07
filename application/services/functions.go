package services

import (
	"github.com/toolfordev/local-cli/application/models"
	"github.com/toolfordev/local-cli/infrastructure/files"
)

func Init() (err error) {
	err = startConfigFile(toolForDevConfig)
	return
}

func Start(filePath string) (err error) {
	configFile := models.ToolForDevFileConfig{}
	err = files.YamlFromFilePathToModel(filePath, &configFile)
	if err != nil {
		return
	}
	err = startConfigFile(configFile)
	return
}

func Destroy() (err error) {
	err = stopConfigFile(toolForDevConfig)
	return
}
