package config

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type conf struct {
	Nakama struct {
		Host   string `yaml:"host"`
		Key    string `yaml:"key"`
		Port   int64  `yaml:"port"`
		Secure bool   `yaml:"secure"`
	}
	Website struct {
		Host string `yaml:"host"`
	}
}

func (c *conf) getConf(env string) *conf {
	yamlFile, err := ioutil.ReadFile(fmt.Sprintf("../build/profiles/%s.yml", env))
	if err != nil {
		log.Printf("yamlFile Get err: #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal err: %v", err)
	}
	fmt.Printf("--- m:\n%v\n\n", c)

	return c
}

// GameConfig go brr
func GameConfig(environment string) {
	fmt.Println("build game config")

	var c conf

	c.getConf(environment)
}
