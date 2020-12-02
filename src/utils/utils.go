package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"types"
)

//FileExists - checks if current file exists
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

//CreateConfigFile - creates a default config file
func CreateConfigFile(fileLocation string) bool {

	//Default config
	config := types.Config{
		ServerPort: "9000",
		Deployments: []types.DeploymentConfig{
			types.DeploymentConfig{
				RepoName:           "test",
				Type:               "repository",
				RepoLocation:       "/REPOSITORY_LOCATION",
				RepoBranch:         "main",
				Secret:             "",
				DeploymentCommands: []types.Command{},
			},
		},
	}

	//Make default config a json object
	data, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		fmt.Println(err)
		return false
	}

	//Create all parent folder if they do not exist.
	pathLocation := strings.Replace(fileLocation, "/config.json", "", 1)
	err = os.MkdirAll(pathLocation, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return false
	}

	//Create the config file.
	err = ioutil.WriteFile(fileLocation, data, 0644)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
