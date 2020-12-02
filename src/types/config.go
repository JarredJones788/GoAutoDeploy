package types

//Command - Command struct
type Command struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}

//DeploymentConfig - Config for a single deployment
type DeploymentConfig struct {
	RepoName           string    `json:"repoName"`
	Type               string    `json:"type"`
	RepoLocation       string    `json:"repoLocation"`
	RepoBranch         string    `json:"repoBranch"`
	Secret             string    `json:"secret"`
	DeploymentCommands []Command `json:"commands"`
	SSHURL             string
}

//Config - start up config
type Config struct {
	ServerPort  string
	Deployments []DeploymentConfig
}
