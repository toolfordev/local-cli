package models

type ToolForDevFileConfig struct {
	ToolForDev ToolForDevConfig `json:"toolForDev,omitempty" yaml:"toolForDev,omitempty"`
}

type ToolForDevConfig struct {
	Project ProjectConfig `json:"project,omitempty" yaml:"project,omitempty"`
}

type ProjectConfig struct {
	Name         string              `json:"name,omitempty" yaml:"name,omitempty"`
	Description  string              `json:"description,omitempty" yaml:"description,omitempty"`
	Variables    []VariableConfig    `json:"variables,omitempty" yaml:"variables,omitempty"`
	Applications []ApplicationConfig `json:"applications,omitempty" yaml:"applications,omitempty"`
	NetworkName  string              `json:"networkName,omitempty" yaml:"networkName,omitempty"`
}

type ApplicationConfig struct {
	Name                 string           `json:"name,omitempty" yaml:"name,omitempty"`
	Description          string           `json:"description,omitempty" yaml:"description,omitempty"`
	ContainerFile        string           `json:"containerFile,omitempty" yaml:"containerFile,omitempty"`
	ContainerName        string           `json:"containerName,omitempty" yaml:"containerName,omitempty"`
	BuildContext         string           `json:"buildContext,omitempty" yaml:"buildContext,omitempty"`
	Ports                []PortConfig     `json:"ports,omitempty" yaml:"ports,omitempty"`
	ExternalPort         string           `json:"externalPort,omitempty" yaml:"externalPort,omitempty"`
	Image                string           `json:"image,omitempty" yaml:"image,omitempty"`
	EnvironmentVariables []VariableConfig `json:"environmentVariables,omitempty" yaml:"environmentVariables,omitempty"`
}

type VariableConfig struct {
	Source     string `json:"source,omitempty" yaml:"source,omitempty"`
	SourceName string `json:"sourceName,omitempty" yaml:"sourceName,omitempty"`
	Name       string `json:"name,omitempty" yaml:"name,omitempty"`
	Value      string `json:"value,omitempty" yaml:"value,omitempty"`
}

type PortConfig struct {
	Application string `json:"application,omitempty" yaml:"application,omitempty"`
	External    string `json:"external,omitempty" yaml:"external,omitempty"`
	Protocol    string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
}
