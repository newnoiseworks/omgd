package utils

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// ProfileConf docs
type ProfileConf struct {
	Nakama struct {
		Host   string `yaml:"host"`
		Key    string `yaml:"key"`
		Port   int    `yaml:"port"`
		Secure bool   `yaml:"secure"`
	}
	Gcloud struct {
		Project string `yaml:"project"`
		Region  string `yaml:"region"`
		Zone    string `yaml:"zone"`
	}
	Firebase struct {
		Project string `yaml:"project"`
	}
	Game struct {
		RealWorldSecondsPerDay string `yaml:"real_world_seconds_per_day"`
		Version                string `yaml:"version"`
	}
	Git struct {
		GameBranch   string `yaml:"game"`
		ServerBranch string `yaml:"server"`
	}
}

// GetProfile d
func GetProfile(env string) ProfileConf {
	c := ProfileConf{}

	yamlFile, err := ioutil.ReadFile(fmt.Sprintf("build/profiles/%s.yml", env))
	if err != nil {
		log.Printf("yamlFile Get err: #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal err: %v", err)
	}

	if env != "local" {
		var key = fmt.Sprintf("the-promised-land-%s-v%s", env, c.Game.Version)
		data := []byte(key)
		c.Nakama.Key = fmt.Sprintf("%x", md5.Sum(data))
	} else {
		c.Nakama.Key = "defaultkey"
	}

	return c
}

// SaveProfile saves that profile to yml
func SaveProfile(profile ProfileConf, env string) {
	yamlBytes, err := yaml.Marshal(&profile)

	if err != nil {
		log.Fatal("Error marshalling from data to saving profile to yaml!")
	}

	err = ioutil.WriteFile(fmt.Sprintf("build/profiles/%s.yml", env), yamlBytes, 0755)

	if err != nil {
		log.Fatal("Error on file write to saving profile to yaml!")
	}
}
