package main

import (
	"deployer"
	"types"
)

func main() {
	config := types.Config{
		ServerPort: ":5645",
		Deployments: []types.DeploymentConfig{
			types.DeploymentConfig{
				Name:     "test",
				Type:     "repository",
				Location: "C:/Users/Jarred/Desktop/testdir",
				Secret:   "123",
				Commands: []types.Command{
					types.Command{
						Name: "pm2",
						Args: []string{"restart", "server"},
					},
				},
			},
		},
	}

	deployer.AutoDeploy{}.Init(&config)
}
