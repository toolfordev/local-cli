package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/toolfordev/local-cli/application/models"
	"github.com/toolfordev/local-cli/infrastructure/docker"
)

var toolForDevConfig models.ToolForDevFileConfig = models.ToolForDevFileConfig{
	ToolForDev: models.ToolForDevConfig{
		Project: models.ProjectConfig{
			Name:        "toolfordev",
			Description: "a toolfordev project",
			Variables: []models.VariableConfig{
				{
					Name:  "EXTENDEDDB_TOOLFORDEV_TYPE",
					Value: "postgres",
				},
				{
					Name:       "EXTENDEDDB_TOOLFORDEV_HOST",
					Source:     "applicationHostname",
					SourceName: "local-database-postgres",
				},
				{
					Name:  "EXTENDEDDB_TOOLFORDEV_PORT",
					Value: "5432",
				},
				{
					Name:  "EXTENDEDDB_TOOLFORDEV_SSL_MODE",
					Value: "prefer",
				},
				{
					Name:  "EXTENDEDDB_TOOLFORDEV_NAME",
					Value: "toolfordev",
				},
				{
					Name:  "EXTENDEDDB_TOOLFORDEV_USER",
					Value: "toolfordev",
				},
				{
					Name:  "EXTENDEDDB_TOOLFORDEV_PASSWORD",
					Value: "Yc7cUNQSkNm9pEuz",
				},
			},
			Applications: []models.ApplicationConfig{
				{
					Name:        "local-database-postgres",
					Description: "a toolfordev database",
					Image:       "quay.io/toolfordev/local-database-postgres:v1.0.3",
					Ports: []models.PortConfig{
						{
							Application: "5432",
							External:    "14432",
						},
					},
					EnvironmentVariables: []models.VariableConfig{
						{
							Name:  "POSTGRES_PASSWORD",
							Value: "NgbxvP9h823wFye8",
						},
						{
							Name:       "TOOLFORDEV_DB",
							Source:     "project",
							SourceName: "EXTENDEDDB_TOOLFORDEV_NAME",
						},
						{
							Name:       "TOOLFORDEV_USER",
							Source:     "project",
							SourceName: "EXTENDEDDB_TOOLFORDEV_USER",
						},
						{
							Name:       "TOOLFORDEV_PASSWORD",
							Source:     "project",
							SourceName: "EXTENDEDDB_TOOLFORDEV_PASSWORD",
						},
					},
				},
				{
					Name:        "local-api-global-variables",
					Description: "a toolfordev api for global variables",
					Image:       "quay.io/toolfordev/local-api-global-variables:v1.0.11",
					Ports: []models.PortConfig{
						{
							Application: "80",
							External:    "14001",
						},
					},
					EnvironmentVariables: []models.VariableConfig{
						{
							Name:   "EXTENDEDDB_TOOLFORDEV_TYPE",
							Source: "project",
						},
						{
							Name:   "EXTENDEDDB_TOOLFORDEV_HOST",
							Source: "project",
						},
						{
							Name:   "EXTENDEDDB_TOOLFORDEV_PORT",
							Source: "project",
						},
						{
							Name:   "EXTENDEDDB_TOOLFORDEV_SSL_MODE",
							Source: "project",
						},
						{
							Name:   "EXTENDEDDB_TOOLFORDEV_NAME",
							Source: "project",
						},
						{
							Name:   "EXTENDEDDB_TOOLFORDEV_USER",
							Source: "project",
						},
						{
							Name:   "EXTENDEDDB_TOOLFORDEV_PASSWORD",
							Source: "project",
						},
					},
				},
			},
		},
	},
}

func startConfigFile(configFile models.ToolForDevFileConfig) (err error) {
	stopConfigFile(configFile)
	loadConfigFileAdvanced(&configFile)
	dockerClient := docker.Docker{}
	err = dockerClient.Init()
	if err != nil {
		return
	}
	err = dockerClient.NetworkCreate(configFile.ToolForDev.Project.NetworkName)
	if err != nil {
		if strings.Contains(err.Error(), fmt.Sprintf("network with name %v already exists", configFile.ToolForDev.Project.NetworkName)) {
			err = nil
		} else {
			return
		}
	}
	for _, application := range configFile.ToolForDev.Project.Applications {
		if application.ContainerFile == "" {
			err = dockerClient.ImagePull(application.Image)
			if err != nil {
				return
			}

			for {
				isReady, err := dockerClient.ImageIsReady(application.Image)
				if err != nil {
					break
				}
				if isReady {
					break
				}
				time.Sleep(time.Second)
			}
			if err != nil {
				return
			}
		}
		err = dockerClient.ContainerCreate(application)
		if err != nil {
			return
		}
		err = dockerClient.ContainerStart(application)
		if err != nil {
			return
		}
		err = dockerClient.NetworkConnect(configFile.ToolForDev.Project.NetworkName, application.ContainerName)
		if err != nil {
			return
		}
	}
	return
}

func stopConfigFile(configFile models.ToolForDevFileConfig) (err error) {
	loadConfigFile(&configFile)
	dockerClient := docker.Docker{}
	err = dockerClient.Init()
	if err != nil {
		return
	}
	for _, application := range configFile.ToolForDev.Project.Applications {
		err = dockerClient.ContainerRemove(application)
		if err != nil {
			if strings.Contains(err.Error(), fmt.Sprintf("No such container: %v", application.ContainerName)) {
				err = nil
			} else {
				return
			}
		}
	}

	err = dockerClient.NetworkRemove(configFile.ToolForDev.Project.NetworkName)
	if err != nil {
		if strings.Contains(err.Error(), fmt.Sprintf("No such network: %v", configFile.ToolForDev.Project.NetworkName)) {
			err = nil
		} else {
			return
		}
	}

	return
}

func loadConfigFile(configFile *models.ToolForDevFileConfig) {
	if configFile.ToolForDev.Project.NetworkName == "" {
		configFile.ToolForDev.Project.NetworkName = fmt.Sprintf("tfd_%v", configFile.ToolForDev.Project.Name)
	}
	for i := range configFile.ToolForDev.Project.Applications {
		if configFile.ToolForDev.Project.Applications[i].ContainerName == "" {
			configFile.ToolForDev.Project.Applications[i].ContainerName = fmt.Sprintf("tfd_%v_%v", configFile.ToolForDev.Project.Name, configFile.ToolForDev.Project.Applications[i].Name)
		}
	}
}

func loadConfigFileAdvanced(configFile *models.ToolForDevFileConfig) {
	loadConfigFile(configFile)
	findVariableValueByName := func(name string) (value string) {
		for i := range configFile.ToolForDev.Project.Variables {
			if configFile.ToolForDev.Project.Variables[i].Name == name {
				value = configFile.ToolForDev.Project.Variables[i].Value
				return
			}
		}
		return
	}
	findApplicationContainerNameByName := func(name string) (containerName string) {
		for i := range configFile.ToolForDev.Project.Applications {
			if configFile.ToolForDev.Project.Applications[i].Name == name {
				containerName = configFile.ToolForDev.Project.Applications[i].ContainerName
				return
			}
		}
		return
	}
	getVarValue := func(source, sourceName string) (value string) {
		switch source {
		case "project":
			value = findVariableValueByName(sourceName)
			return
		case "applicationHostname":
			value = fmt.Sprintf("%v.local.toolfor.dev", findApplicationContainerNameByName(sourceName))
			return
		}
		return
	}
	for i := range configFile.ToolForDev.Project.Variables {
		if configFile.ToolForDev.Project.Variables[i].Source == "" {
			continue
		}

		if configFile.ToolForDev.Project.Variables[i].SourceName == "" {
			configFile.ToolForDev.Project.Variables[i].SourceName = configFile.ToolForDev.Project.Variables[i].Name
		}

		configFile.ToolForDev.Project.Variables[i].Value = getVarValue(configFile.ToolForDev.Project.Variables[i].Source, configFile.ToolForDev.Project.Variables[i].SourceName)
	}

	for i := range configFile.ToolForDev.Project.Applications {
		for j := range configFile.ToolForDev.Project.Applications[i].Ports {
			if configFile.ToolForDev.Project.Applications[i].Ports[j].External == "" {
				configFile.ToolForDev.Project.Applications[i].Ports[j].External = configFile.ToolForDev.Project.Applications[i].Ports[j].Application
			}
			if configFile.ToolForDev.Project.Applications[i].Ports[j].Protocol == "" {
				configFile.ToolForDev.Project.Applications[i].Ports[j].Protocol = "tcp"
			}
		}
		for j := range configFile.ToolForDev.Project.Applications[i].EnvironmentVariables {
			if configFile.ToolForDev.Project.Applications[i].EnvironmentVariables[j].Source == "" {
				continue
			}

			if configFile.ToolForDev.Project.Applications[i].EnvironmentVariables[j].SourceName == "" {
				configFile.ToolForDev.Project.Applications[i].EnvironmentVariables[j].SourceName = configFile.ToolForDev.Project.Applications[i].EnvironmentVariables[j].Name
			}

			configFile.ToolForDev.Project.Applications[i].EnvironmentVariables[j].Value = getVarValue(configFile.ToolForDev.Project.Applications[i].EnvironmentVariables[j].Source, configFile.ToolForDev.Project.Applications[i].EnvironmentVariables[j].SourceName)
		}
	}
}
