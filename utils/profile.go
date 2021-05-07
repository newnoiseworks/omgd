package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
)

type Command struct {
	Cmd string `yaml:"cmd"`
	Dir string `yaml:"dir"`
}

type CommandConfig struct {
	Name  string    `yaml:"name"`
	Dir   string    `yaml:"dir"`
	Steps []Command `yaml:"steps"`
}

type ProfileConf struct {
	Name string
	Git  struct {
		GameBranch string `yaml:"branch"`
		Repo       string `yaml:"repo"`
	}
	Main  []CommandConfig `yaml:"main"`
	Tasks []CommandConfig `yaml:"tasks"`
}

func GetProfileAsMap(env string) *map[interface{}]interface{} {
	profile := GetProfile(env)

	yamlBytes, err := yaml.Marshal(&profile)

	if err != nil {
		log.Fatal("Error marshalling from data back to yaml!")
	}

	c := make(map[interface{}]interface{})

	err = yaml.Unmarshal(yamlBytes, &c)
	if err != nil {
		log.Fatalf("Unmarshal err: %v", err)
	}

	splits := strings.Split(env, "/")
	c["name"] = splits[len(splits)-1]

	return &c
}

func GetProfile(env string) *ProfileConf {
	c := ProfileConf{}

	yamlFile, err := ioutil.ReadFile(fmt.Sprintf("%s.yml", env))
	if err != nil {
		log.Printf("yamlFile Get err: #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal err: %v", err)
	}

	splits := strings.Split(env, "/")
	c.Name = splits[len(splits)-1]

	return &c
}

// SaveProfile saves that profile to yml
func SaveProfile(profile ProfileConf, env string) {
	yamlBytes, err := yaml.Marshal(&profile)

	if err != nil {
		log.Fatal("Error marshalling from data to saving profile to yaml!")
	}

	err = ioutil.WriteFile(fmt.Sprintf("profiles/%s.yml", env), yamlBytes, 0755)

	if err != nil {
		log.Fatal("Error on file write to saving profile to yaml!")
	}
}
