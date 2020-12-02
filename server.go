package main

import (
	"deployer"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime"
	"types"
	"utils"
)

func main() {

	//Window config location
	configLocation := "./cfg/config.json"

	//Linux config location
	if runtime.GOOS == "linux" {
		configLocation = "/etc/GoAutoDeploy/config.json"
	}

	//If a config file does not exists. Create a default one.
	if !utils.FileExists(configLocation) {
		if utils.CreateConfigFile(configLocation) {
			fmt.Println("Default config created. You can find it here: " + configLocation)
		} else {
			fmt.Println("Failed creating default configuration file..")
			return
		}
	}

	//Read config file
	configFile, err := ioutil.ReadFile(configLocation)
	if err != nil {
		fmt.Println(err)
		return
	}

	var config types.Config

	//Cast config into types.Config struct
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		fmt.Println("Failed reading config file")
		fmt.Println(err)
		return
	}

	//Start AutoDeploy service
	err = deployer.AutoDeploy{}.Init(&config)
	if err != nil {
		fmt.Println(err.Error())
	}
}
