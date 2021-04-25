package utils

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// MissionItem struct
type MissionItem struct {
	Key      string `yaml:"key"`
	Quantity int    `yaml:"quantity"`
}

// Mission struct
type Mission struct {
	Key     string        `yaml:"key"`
	Title   string        `yaml:"title"`
	Prereqs string        `yaml:"prereqs"`
	Reqs    []MissionItem `yaml:"reqs"`
	Awards  []MissionItem `yaml:"awards"`
}

// MissionData dem missions
type MissionData struct {
	Missions []Mission `yaml:"missions"`
}

// GetMissions d
func GetMissions() MissionData {
	c := MissionData{}

	yamlFile, err := ioutil.ReadFile("resources/missions.yml")
	if err != nil {
		log.Printf("yamlFile Get err: #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal err: %v", err)
	}

	return c
}