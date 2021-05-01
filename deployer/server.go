package deployer

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/newnoiseworks/tpl-fred/utils"
)

type Server struct {
	Environment string
	OutputDir   string
	VolumeReset bool
	CmdOnDir    func(string, string, string)
}

func (ds Server) Deploy() {
	fmt.Println("deploying server")

	serverPath, err := filepath.Abs(fmt.Sprintf("%s/server/deploy/gcp", ds.OutputDir))
	if err != nil {
		log.Fatal(err)
		return
	}

	switch ds.Environment {
	case "local":
		fmt.Println("Need to make local deployment commands")
		break
	case "production":
		fmt.Println("Do special production deployment stuff?")
		ds.deployServerBasedOnProfile(serverPath)
		break
	default:
		ds.deployServerBasedOnProfile(serverPath)
		break
	}
}

func (ds Server) deployServerBasedOnProfile(serverPath string) {
	var config = utils.GetProfile(ds.Environment)

	var cmdString = fmt.Sprintf("GCP_UPDATE=true GCP_PROJECT=%s GCP_ZONE=%s ./deploy.sh", config.Gcloud.Project, config.Gcloud.Zone)

	if ds.VolumeReset {
		cmdString = "GCP_UPDATE_REMOVE_VOLUME=true " + cmdString
	}

	ds.CmdOnDir(
		cmdString,
		"running deploy.sh script in server dir",
		serverPath,
	)
}
