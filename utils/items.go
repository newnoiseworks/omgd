package utils

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// GetItems d
func GetItems() map[interface{}]interface{} {
	c := make(map[interface{}]interface{})

	yamlFile, err := ioutil.ReadFile("resources/items.yml")
	if err != nil {
		log.Printf("yamlFile Get err: #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal err: %v", err)
	}

	return c
}
