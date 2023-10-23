package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime/debug"
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
	Name    string
	Main    []CommandConfig `yaml:"main"`
	Tasks   []CommandConfig `yaml:"tasks"`
	path    string
	rootDir string
}

func (pc ProfileConf) getTopLevelOMGDProfileAsMap() map[interface{}]interface{} {
	c := make(map[interface{}]interface{})

	omgdProfilePath := strings.Replace(pc.path, pc.Name, "omgd", 1)

	if pc.rootDir != "" {
		omgdProfilePath = filepath.Join(pc.rootDir, omgdProfilePath)
	}

	yamlFile, err := ioutil.ReadFile(omgdProfilePath)

	if err != nil {
		return c
	}
	err = yaml.Unmarshal(yamlFile, &c)

	if err != nil {
		log.Fatalf("Unmarshal err: %v", err)
	}

	return c
}

func (pc ProfileConf) getRootProfileAsMap() map[interface{}]interface{} {
	c := make(map[interface{}]interface{})

	profilePath := pc.path

	if pc.rootDir != "" {
		profilePath = filepath.Join(pc.rootDir, profilePath)
	}

	yamlFile, err := ioutil.ReadFile(profilePath)

	if err != nil {
		log.Fatalf("could not read file: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &c)

	if err != nil {
		log.Fatalf("Unmarshal err: %v", err)
	}

	return c
}

func (pc ProfileConf) GetProfileAsMap() map[interface{}]interface{} {
	path := pc.path

	if pc.rootDir != "" {
		path = filepath.Join(pc.rootDir, path)
	}

	omgdFile := pc.getTopLevelOMGDProfileAsMap()

	c := pc.getRootProfileAsMap()

	mergerecursive(&c, &omgdFile, 4)

	return omgdFile
}

func GetProfile(path string) *ProfileConf {
	c := ProfileConf{
		path: path,
	}

	yamlFile, err := os.ReadFile(c.path)
	if err != nil {
		debug.PrintStack()
		log.Fatalf("yamlFile Get err: #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal err: %v", err)
	}

	splits := strings.Split(path, "/")
	c.Name = strings.Replace(splits[len(splits)-1], ".yml", "", 1)

	return &c
}

func GetProfileFromDir(path string, dir string) *ProfileConf {
	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	err = os.Chdir(filepath.Join(root, dir))
	if err != nil {
		log.Fatal(err)
	}

	profile := GetProfile(path)
	profile.rootDir = dir

	err = os.Chdir(root)
	if err != nil {
		log.Fatal(err)
	}

	return profile
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
	profileMap := profile.getRootProfileAsMap()
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
			profile := GetProfile(fmt.Sprintf("%s/profiles/%s", dir, file.Name()))

			BuildTemplateFromPath(
				fmt.Sprintf("%s/profiles/profile.yml.omgdptpl", dir),
				profile,
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

// c/o https://github.com/wiebew/golang_merge_yml/blob/main/merge_yaml.go
func isMap(typeName string) bool {
	switch typeName {
	case "map[interface {}]interface {}":
		return true
	default:
		return false
	}
}

// c/o https://github.com/wiebew/golang_merge_yml/blob/main/merge_yaml.go
func mergerecursive(master *map[interface{}]interface{}, merge *map[interface{}]interface{}, level int) {

	for k, v := range *master {
		_, exists := (*merge)[k]
		if exists {
			// key exist in the target yaml
			// if it is a map we need to (recursively) descend into it to check every value in the map
			// this prevents losing values if the master only has one underlying value and the default multiple
			if isMap(reflect.TypeOf(v).String()) {
				masternode := v.(map[interface{}]interface{}) // type assertion (typecast)
				// check if both types are a map types
				if isMap(reflect.TypeOf((*merge)[k]).String()) {
					mergenode := (*merge)[k].(map[interface{}]interface{}) // type assertion
					mergerecursive(&masternode, &mergenode, level+1)
				} else {
					log.Fatal("Key [", k, "] is map/list of values in one yaml and a singular value in the other yaml, can't merge them")
				}
			} else {
				// key is not a map, so we just need to copy the value if they are both non-map types
				if !isMap(reflect.TypeOf((*merge)[k]).String()) {
					(*merge)[k] = v
				} else {
					log.Fatal("Key [", k, "] is map/list of values in one yaml and a singular value in the other yaml, can't merge them")
				}
			}
		} else {
			// key does not exists in the target, just add the whole node/value
			(*merge)[k] = v
		}
	}
}

// c/o https://github.com/wiebew/golang_merge_yml/blob/main/merge_yaml.go
func merge(masterpath *string, defaultspath *string, merge *map[interface{}]interface{}) {
	// will merge values of both yaml files into merge map
	// values in master will overrule values in defaults

	var master map[interface{}]interface{}
	var defaults map[interface{}]interface{}

	bs, err := ioutil.ReadFile(*masterpath)
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(bs, &master); err != nil {
		panic(err)
	}

	bs, err = ioutil.ReadFile(*defaultspath)

	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(bs, &defaults); err != nil {
		panic(err)
	}
	*merge = defaults
	mergerecursive(&master, merge, 0)
}
