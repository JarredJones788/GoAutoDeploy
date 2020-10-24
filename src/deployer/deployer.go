package deployer

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"types"

	"github.com/gorilla/mux"
)

//AutoDeploy type
type AutoDeploy struct {
	Deployments *[]types.DeploymentConfig
}

//Init - inits all routes.
func (deploy AutoDeploy) Init(config *types.Config) error {

	deploy.Deployments = &config.Deployments

	//Setup mux deploy
	r := mux.NewRouter()
	deploy.setUpRoutes(r)
	fmt.Println("Server Started")
	err := http.ListenAndServe(config.ServerPort, r)
	if err != nil {
		return err
	}

	return nil
}

//setUpRoutes - sets up all endpoints for the service
func (deploy AutoDeploy) setUpRoutes(r *mux.Router) {
	r.HandleFunc("/github/push", deploy.updatePushed)
}

//setUpHeaders - sets the desired headers for an http response
func (deploy AutoDeploy) setUpHeaders(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Max-Age", "120")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	if r.Method == http.MethodOptions {
		w.WriteHeader(200)
		return false
	}
	return true
}

//validateRequest - validates if the requested repository is known in the config.
//validates if the request is from github.
func (deploy AutoDeploy) validateRequest(r *http.Request) (*types.DeploymentConfig, error) {

	//Read request body into a byte array to calculate hash.
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil || len(payload) == 0 {
		return nil, errors.New("Request payload is empty")
	}

	//Read payload into a json object to get repository details
	var requestPayload types.Payload
	if err := json.Unmarshal(payload, &requestPayload); err != nil {
		return nil, errors.New("Failed parsing payload")
	}

	var selectedRepo *types.DeploymentConfig

	//Loop over all repositories provided in the config.
	//Find the matching repo's.
	for _, repo := range *deploy.Deployments {
		if repo.Name == requestPayload.Repository.Name {
			selectedRepo = &repo
			selectedRepo.SSHURL = requestPayload.Repository.SSHURL
			break
		}
	}

	//No repo was found
	if selectedRepo == nil {
		return nil, errors.New("Repository in the request does not match the ones given in config")
	}

	//if a secret is passed then verify the request was sent from github
	if len(selectedRepo.Secret) > 0 {
		signature := r.Header.Get("X-Hub-Signature")
		if len(signature) == 0 {
			return nil, errors.New("No Signature was found on request for repo: " + selectedRepo.Name)
		}
		mac := hmac.New(sha1.New, []byte(selectedRepo.Secret))
		_, _ = mac.Write(payload)
		expectedMAC := hex.EncodeToString(mac.Sum(nil))

		if !hmac.Equal([]byte(signature[5:]), []byte(expectedMAC)) {
			return nil, errors.New("Invalid Signature on request for repo: " + selectedRepo.Name)
		}
	}

	return selectedRepo, nil
}

//updatePushed - called by github web hook when a new update was pushed
func (deploy AutoDeploy) updatePushed(w http.ResponseWriter, r *http.Request) {
	if !deploy.setUpHeaders(w, r) {
		return //request was an OPTIONS which was handled.
	}

	//Check if the request is a from github
	repo, err := deploy.validateRequest(r)
	if err != nil {
		fmt.Println(err.Error())
		w.Write([]byte("Request Failed"))
		return
	}

	//Check that a git repo exists at the location
	cmd := exec.Command("git", "-C", repo.Location, "rev-parse", "--show-toplevel")
	_, err = cmd.Output()
	if err != nil {
		fmt.Println("Not Git Repository was found at: " + repo.Location)
		w.Write([]byte("Request Failed"))
		return
	}

	fmt.Println("Location is a valid Git Repo.")

	//Pull the current repo to the location
	cmd = exec.Command("git", "-C", repo.Location, "pull", repo.SSHURL)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	_, err = cmd.Output()
	if err != nil {
		fmt.Println("Failed pulling repo " + repo.Name + " -> " + stderr.String())
		w.Write([]byte("Request Failed"))
		return
	}

	fmt.Println("Git Repo '" + repo.Name + "' pulled succesfully")

	//Run commands that are provided for the deployment
	for _, command := range repo.Commands {
		cmd := exec.Command(command.Name, command.Args...)
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Command Failed: " + command.Name + " -> " + err.Error())
			continue
		}
		fmt.Println("Command " + command.Name + " was ran:")
		fmt.Println(string(output))
	}

	w.Write([]byte("Success"))
}
