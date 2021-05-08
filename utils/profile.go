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
	path  string
	env   string
}

func (pc ProfileConf) GetProfileAsMap() map[interface{}]interface{} {
	c := make(map[interface{}]interface{})

	yamlFile, err := ioutil.ReadFile(pc.path)

	if err != nil {
		log.Printf("yamlFile Get err: #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &c)

	if err != nil {
		log.Fatalf("Unmarshal err: %v", err)
	}

	return c
}

func GetProfile(env string) *ProfileConf {
	c := ProfileConf{}

	c.env = env
	c.path = fmt.Sprintf("%s.yml", env)

	yamlFile, err := ioutil.ReadFile(c.path)
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

func (profile ProfileConf) SaveProfileFromMap(profileMap *map[interface{}]interface{}) {
	yamlBytes, err := yaml.Marshal(profileMap)

	if err != nil {
		log.Fatal("Error marshalling from data to saving profile to yaml!")
	}

	err = ioutil.WriteFile(profile.path, yamlBytes, 0755)

	if err != nil {
		log.Fatal("Error on file write to saving profile to yaml!")
	}
}
