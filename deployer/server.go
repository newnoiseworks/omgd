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

	var cmd = "down"

	if volumeReset {
		cmd = "down -v"
	}

	runDockerComposeCmdOnServerDir(
		cmd,
		"docker-compose down on remote nakama servers",
		config,
	)

	runCmdOnServerDir(
		"rm -rf ./nakama",
		"rm on remote nakama folder",
		config,
	)

	utils.CmdOnDir(
		fmt.Sprintf(`gcloud compute scp --project %s --zone %s --recurse --force-key-file-overwrite ./nakama %s:`, config.Gcloud.Project, config.Gcloud.CloudZone, config.Gcloud.Instance),
		"scp up nakama module folder",
		serverPath,
	)

	utils.CmdOnDir(
		fmt.Sprintf(`gcloud compute scp --project %s --zone %s --force-key-file-overwrite docker-compose.yml %s:`, config.Gcloud.Project, config.Gcloud.CloudZone, config.Gcloud.Instance),
		"scp up docker-compose file",
		serverPath,
	)

	runDockerComposeCmdOnServerDir(
		`up -d`,
		"docker-compose up on remote nakama servers",
		config,
	)
}

func runDockerComposeCmdOnServerDir(cmdStr string, cmdDesc string, config utils.ProfileConf) {
	utils.CmdOnDir(
		fmt.Sprintf(`gcloud compute ssh --zone %s --project %s --command "docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v "\$PWD:\$PWD" -w="\$PWD" docker/compose:1.24.0 %s" %s`, config.Gcloud.CloudZone, config.Gcloud.Project, cmdStr, config.Gcloud.Instance),
		cmdDesc,
		serverPath,
	)
}

func runCmdOnServerDir(cmdStr string, cmdDesc string, config utils.ProfileConf) {
	utils.CmdOnDir(
		fmt.Sprintf(`gcloud compute ssh --zone %s --project %s --command "%s" %s`, config.Gcloud.CloudZone, config.Gcloud.Project, cmdStr, config.Gcloud.Instance),
		cmdDesc,
		serverPath,
	)
}
