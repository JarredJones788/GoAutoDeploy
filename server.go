package main

import (
	"bytes"
	"deployer"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"types"
	"utils"

	"github.com/kardianos/osext"
)

func main() {

	//If no args are passed then start server normally
	//If args are passed then complete the required task.
	if len(os.Args) <= 1 {
		startServer()
	} else {
		switch os.Args[1] {
		case "init":
			initService()
			break
		case "start":
			startService()
			break
		case "stop":
			stopService()
			break
		case "reload":
		case "restart":
			restartService()
			break
		case "remove":
			removeService()
			break
		case "status":
			serviceStatus()
			break
		default:
			fmt.Println("Invalid base command")
			break
		}
	}
}

//Creates a systemd service
//Only works on linux
func initService() {
	if runtime.GOOS == "windows" {
		fmt.Println("Init only works on linux")
		return
	}

	path, err := osext.ExecutableFolder()
	if err != nil {
		fmt.Println(err)
		return
	}

	content := "[Unit]\nDescription=AutoDeploy\n\n[Service]\nType=simple\nRestart=on-failure\nRestartSec=5s\nUser=root\nExecStart=" + path + "/autodeploy\n\n[Install]\nWantedBy=multi-user.target"
	if ioutil.WriteFile("/etc/systemd/system/autodeploy.service", []byte(content), 0644) != nil {
		fmt.Println("Error saving service file")
		return
	}
	cmd := exec.Command("systemctl", "enable", "autodeploy.service")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	_, err = cmd.Output()
	if err != nil {
		fmt.Println(stderr.String())
		return
	}
	cmd = exec.Command("systemctl", "start", "autodeploy")
	cmd.Stderr = &stderr
	_, err = cmd.Output()
	if err != nil {
		fmt.Println(stderr.String())
		return
	}

	fmt.Println("AutoDeploy service was created")
}

//Restarts systemd service
func restartService() {

	if runtime.GOOS == "windows" {
		fmt.Println("Restart only works on linux")
		return
	}

	cmd := exec.Command("systemctl", "restart", "autodeploy.service")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(stderr.String())
		return
	}

	fmt.Println("Service restarted")
}

//Stops systemd service
func stopService() {

	if runtime.GOOS == "windows" {
		fmt.Println("Stop only works on linux")
		return
	}

	cmd := exec.Command("systemctl", "stop", "autodeploy.service")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(stderr.String())
		return
	}

	fmt.Println("Service stopped")
}

//removes systemd service
func removeService() {
	if runtime.GOOS == "windows" {
		fmt.Println("Remove only works on linux")
		return
	}

	cmd := exec.Command("systemctl", "stop", "autodeploy.service")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(stderr.String())
		return
	}

	cmd = exec.Command("systemctl", "disable", "autodeploy.service")
	cmd.Stderr = &stderr
	_, err = cmd.Output()
	if err != nil {
		fmt.Println(stderr.String())
		return
	}

	cmd = exec.Command("rm", "/etc/systemd/system/autodeploy.service")
	cmd.Stderr = &stderr
	_, err = cmd.Output()
	if err != nil {
		fmt.Println(stderr.String())
		return
	}

	fmt.Println("AutoDeploy service has been removed")
}

//get status of the service
func serviceStatus() {
	if runtime.GOOS == "windows" {
		fmt.Println("Status only works on linux")
		return
	}

	cmd := exec.Command("systemctl", "status", "autodeploy.service")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	status, err := cmd.Output()
	if err != nil {
		fmt.Println(stderr.String())
		return
	}

	fmt.Println(string(status))
}

//starts the service
func startService() {

	if runtime.GOOS == "windows" {
		fmt.Println("Start only works on linux")
		return
	}

	cmd := exec.Command("systemctl", "start", "autodeploy.service")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(stderr.String())
		return
	}

	fmt.Println("Service started")
}

//Start the server without a service.
func startServer() {

	//Window config location
	configLocation := "./cfg/config.json"

	//Linux config location
	if runtime.GOOS == "linux" {
		configLocation = "/etc/autodeploy/config.json"
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
