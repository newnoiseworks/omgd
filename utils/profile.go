package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Command struct {
	Cmd  string `yaml:"cmd"`
	Dir  string `yaml:"dir"`
	Desc string `yaml:"desc"`
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
	c := ProfileConf{
		env:  env,
		path: fmt.Sprintf("%s.yml", env),
	}

	yamlFile, err := os.ReadFile(c.path)
	if err != nil {
		log.Fatalf("yamlFile Get err: #%v ", err)
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
		// log.Fatal("Error on file write to saving profile to yaml!")
		log.Fatal(err)
	}
}

func (profile ProfileConf) Get(key string) interface{} {
	profileMap := profile.GetProfileAsMap()
	keys := strings.Split(key, ".")
	return getValueToKeyWithArray(keys, 0, profileMap)
}

func (profile ProfileConf) GetArray(key string) []interface{} {
	profileMap := profile.GetProfileAsMap()
	keys := strings.Split(key, ".")
	return getValueToKeyWithArray(keys, 0, profileMap).([]interface{})
}

func (profile ProfileConf) UpdateProfile(key string, val interface{}) {
	profileMap := profile.GetProfileAsMap()
	keys := strings.Split(key, ".")
	setValueToKeyWithArray(keys, 0, profileMap, val)
	profile.SaveProfileFromMap(&profileMap)
}

func setValueToKeyWithArray(keys []string, keyIndex int, obj map[interface{}]interface{}, value interface{}) {
	for k, v := range obj {
		if key, ok := k.(string); ok {
			if key == keys[keyIndex] {
				if keyIndex == len(keys)-1 {
					obj[k] = value
				} else {
					setValueToKeyWithArray(keys, keyIndex+1, v.(map[interface{}]interface{}), value)
				}
			}
		}
	}
}

func getValueToKeyWithArray(keys []string, keyIndex int, obj map[interface{}]interface{}) interface{} {
	for k, v := range obj {
		if key, ok := k.(string); ok {
			if key == keys[keyIndex] {
				if keyIndex == len(keys)-1 {
					return obj[k]
				} else {
					return getValueToKeyWithArray(keys, keyIndex+1, v.(map[interface{}]interface{}))
				}
			}
		}
	}

	return nil
}

func BuildProfiles(dir string, verbose bool) {
	files, err := ioutil.ReadDir(fmt.Sprintf("%s/profiles", dir))
	if err != nil {
		log.Fatal(err)
	}

	err = os.Mkdir(fmt.Sprintf("%s/.omgd", dir), 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	for _, file := range files {
		splits := strings.Split(file.Name(), ".")
		ext := splits[len(splits)-1]

		if ext == "yml" && splits[0] != "example" {
			BuildTemplateFromPath(
				fmt.Sprintf("%s/profiles/profile.yml.omgdptpl", dir),
				fmt.Sprintf("%s/profiles/%s", dir, strings.Replace(file.Name(), ".yml", "", 1)),
				fmt.Sprintf("%s/profiles", dir),
				"omgdptpl",
				false,
				verbose,
			)

			oldPath := fmt.Sprintf("%s/profiles/profile.yml", dir)
			newPath := fmt.Sprintf("%s/.omgd/%s", dir, file.Name())

			if verbose {
				log.Printf("moving file from %s >> %s", oldPath, newPath)
			}

			origFile, err := os.ReadFile(oldPath)

			if err != nil {
				log.Fatal(err)
			}

			err = os.WriteFile(newPath, origFile, 0755)

			if err != nil {
				log.Fatal(err)
			}

			err = os.Remove(oldPath)

			if err != nil {
				log.Fatal(err)
			}

			_, err = os.Stat(newPath)

			if err != nil {
				log.Fatal(err)
			} else if verbose {
				log.Println("file successfully moved")
			}
		}
	}
}
