package deployer

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/newnoiseworks/tpl-fred/utils"
)

var serverPath string

// DeployServer d
func DeployServer(environment string, buildPath string, volumeReset bool) {
	fmt.Println("deploying server")

	_serverPath, err := filepath.Abs(fmt.Sprintf("%s/server", buildPath))
	if err != nil {
		log.Fatal(err)
		return
	}

	serverPath = _serverPath

	switch environment {
	case "local":
		fmt.Println("Need to make local deployment commands")
		break
	case "production":
		fmt.Println("Do special production deployment stuff?")
		deployServerBasedOnProfile(environment, buildPath, volumeReset)
		break
	default:
		deployServerBasedOnProfile(environment, buildPath, volumeReset)
		break
	}
}

func deployServerBasedOnProfile(environment string, buildPath string, volumeReset bool) {
	var config = utils.GetProfile(environment)

	var cmdString = fmt.Sprintf("GCP_UPDATE=true GCP_PROJECT=%s GCP_ZONE=%s ./gcp_init.sh", config.Gcloud.Project, config.Gcloud.Zone)

	if volumeReset {
		cmdString = "GCP_UPDATE_REMOVE_VOLUME=true " + cmdString
	}

	utils.CmdOnDir(
		cmdString,
		"running gcp_init.sh script in server dir",
		serverPath,
	)
}
