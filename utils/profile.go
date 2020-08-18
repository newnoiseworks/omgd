package utils

import (
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
	Website struct {
		Host string `yaml:"host"`
	}
	Game struct {
		RealWorldSecondsPerDay string `yaml:"real_world_seconds_per_day"`
		Version                string `yaml:"version"`
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

	return c
}
