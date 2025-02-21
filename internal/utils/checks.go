package utils

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/NethermindEth/1click/configs"
	log "github.com/sirupsen/logrus"
)

/*
CheckDependencies :
This function is responsible for checking if on-premise setup dependencies are installed on host machine

params :-
a. dependencies []string
List of dependencies to be checked

returns :-
a. []string
List of dependencies that are not installed
*/
func CheckDependencies(dependencies []string) (pending []string) {
	for _, dependency := range dependencies {
		_, err := exec.LookPath(dependency)
		if err != nil {
			log.Errorf(configs.DependencyNotInstalledError, dependency)
			pending = append(pending, dependency)
		}
	}
	return
}

/*
PreCheck :
Check if docker-compose can be used to interact with the generated docker-compose script

params :-
a. generationPath string
Path to the generated docker-compose script

returns :-
a. error
Error if any
*/
func PreCheck(generationPath string) error {
	// Check that docker and docker-compose are installed
	pending := CheckDependencies([]string{"docker", "docker-compose"})
	for _, dependency := range pending {
		log.Errorf(configs.DependencyNotInstalledError, dependency)
	}
	if len(pending) > 0 {
		return fmt.Errorf(configs.DependenciesMissingError)
	}

	// Check docker engine is on
	log.Debugf(configs.RunningCommand, configs.DockerPsCMD)
	if _, err := RunCmd(configs.DockerPsCMD, true, false); err != nil {
		return fmt.Errorf(configs.DockerEngineOffError, err)
	}

	// Check if docker-compose script was generated
	file := generationPath + "/" + configs.DefaultDockerComposeScriptName
	if _, err := os.Stat(file); os.IsNotExist(err) {
		log.Errorf(configs.OpeningFileError, file, err)
		return fmt.Errorf(configs.DockerComposeScriptNotFoundError, generationPath, configs.DefaultDockerComposeScriptsPath)
	}

	return nil
}

/*
CheckContainers :
Check if containers of generated docker-compose script are running

params :-
a. generationPath string
Path to the generated docker-compose script

returns :-
a. string
Output of 'docker ps --services --filter status=running'
b. error
Error if any
*/
func CheckContainers(generationPath string) (string, error) {
	// Check if docker-compose script is running
	psCMD := fmt.Sprintf(configs.DockerComposePsServicesCMD, generationPath+"/"+configs.DefaultDockerComposeScriptName)
	log.Debugf(configs.RunningCommand, psCMD)
	rawServices, err := RunCmd(psCMD, true, false)
	if err != nil || rawServices == "\n" {
		if rawServices == "\n" && err == nil {
			err = fmt.Errorf(configs.DockerComposePsReturnedEmptyError)
		}
		return "", fmt.Errorf(configs.ScriptIsNotRunningError, err)
	}

	return rawServices, nil
}
