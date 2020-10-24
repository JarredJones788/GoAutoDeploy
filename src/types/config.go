package types

//Command - Command struct
type Command struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}

//DeploymentConfig - Config for a single deployment
type DeploymentConfig struct {
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Location string    `json:"location"`
	Secret   string    `json:"secret"`
	Commands []Command `json:"commands"`
	SSHURL   string
}

//Config - start up config
type Config struct {
	ServerPort  string
	Deployments []DeploymentConfig
}
