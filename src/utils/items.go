package utils

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type priceItem struct {
	Item     string `yaml:"item"`
	Quantity int    `yaml:"quantity"`
}

type item struct {
	Key                string      `yaml:"key"`
	Sale               []priceItem `yaml:"sale"`
	Price              []priceItem `yaml:"price"`
	PurchaseValidation bool        `yaml:"purchase_validation"`
	ComboPrice         bool        `yaml:"combo_price"`
	ToPlant            string      `yaml:"toPlant"`
	FromSeeds          string      `yaml:"fromSeeds"`
	GrowthStages       []int       `yaml:"growthStages"`
	MaxYield           int         `yaml:"maxYield"`
}

// InventoryItems docs
type InventoryItems struct {
	Items []item `yaml:"items"`
}

// GetItems d
func GetItems() InventoryItems {
	c := InventoryItems{}

	yamlFile, err := ioutil.ReadFile("../resources/items.yml")
	if err != nil {
		log.Printf("yamlFile Get err: #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal err: %v", err)
	}

	return c
}
