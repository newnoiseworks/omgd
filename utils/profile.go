package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"gopkg.in/yaml.v2"
)

type ProfileConf struct {
	Name        string
	OMGD        OMGDConfig `yaml:"omgd"`
	OMGDProfile *ProfileConf
	path        string
	rootDir     string
}

type OMGDConfig struct {
	Name    string        `yaml:"name"`
	Game    GameConfig    `yaml:"game"`
	Servers ServersConfig `yaml:"servers"`
	GCP     GCPConfig     `yaml:"gcp"`
}

type ServersConfig struct {
	Services []ServerService `yaml:"services"`
	Ports    PortConfig      `yaml:"ports"`
	Host     string          `yaml:"host"`
}

type ServerService struct {
	BuildService string `yaml:"build-service"`
}

type PortConfig struct {
	TCP string `yaml:"tcp"`
	UDP string `yaml:"udp"`
}

type GameConfig struct {
	Targets []GameTargetConfig `yaml:"targets"`
}

type GameTargetConfig struct {
	BuildService string `yaml:"build-service"`
	Copy         string `yaml:"copy"`
	To           string `yaml:"to"`
}

type GCPConfig struct {
	Project   string `yaml:"project"`
	Zone      string `yaml:"zone"`
	Bucket    string `yaml:"bucket"`
	CredsFile string `yaml:"creds-file"`
}

func (pc ProfileConf) getProfileAsMapFromPath(profilePath string) (map[interface{}]interface{}, error) {
	c := make(map[interface{}]interface{})

	if pc.rootDir != "" {
		profilePath = filepath.Join(pc.rootDir, profilePath)
	}

	yamlFile, err := ioutil.ReadFile(profilePath)

	if err != nil {
		return c, err
	}

	err = yaml.Unmarshal(yamlFile, &c)

	if err != nil {
		return c, err
	}

	return c, nil
}

func (pc ProfileConf) getTopLevelOMGDProfileAsMap() map[interface{}]interface{} {
	omgdProfilePath := strings.Replace(pc.path, pc.Name, "omgd", 1)

	conf, err := pc.getProfileAsMapFromPath(omgdProfilePath)

	if err != nil {
		LogTrace(fmt.Sprintf("Error getting top level OMGD profile as map, %v", err))
	}

	return conf
}

func (pc ProfileConf) getOMGDCloudProfileAsMap() map[interface{}]interface{} {
	omgdCloudProfilePath := strings.Replace(pc.path, pc.Name, "omgd.cloud", 1)

	conf, err := pc.getProfileAsMapFromPath(omgdCloudProfilePath)

	if err != nil {
		LogTrace(fmt.Sprintf("Error getting OMGD cloud profile as map, %v", err))
	}

	return conf
}

func (pc ProfileConf) getRootProfileAsMap() map[interface{}]interface{} {
	profilePath := pc.path

	conf, err := pc.getProfileAsMapFromPath(profilePath)

	if err != nil {
		LogFatal(fmt.Sprintf("Error getting profile as map, %v", err))
	}

	return conf
}

func (pc ProfileConf) GetProfileAsMap() map[interface{}]interface{} {
	path := pc.path

	if pc.rootDir != "" {
		path = filepath.Join(pc.rootDir, path)
	}

	omgdFile := pc.getTopLevelOMGDProfileAsMap()
	cloudFile := pc.getOMGDCloudProfileAsMap()
	profile := pc.getRootProfileAsMap()

	mergerecursive(&cloudFile, &profile, 4)

	mergerecursive(&profile, &omgdFile, 4)

	omgdFile["Name"] = pc.Name

	return omgdFile
}

func GetProfile(path string) *ProfileConf {
	c := ProfileConf{
		path: path,
	}

	splits := strings.Split(path, string(os.PathSeparator))
	c.Name = strings.Replace(splits[len(splits)-1], ".yml", "", 1)

	var bytes, err = yaml.Marshal(c.GetProfileAsMap())
	if err != nil {
		LogFatal(fmt.Sprintf("YAML marshal err from map: %v", err))
	}

	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		LogFatal(fmt.Sprintf("YAML Unmarshal err: %v", err))
	}

	if c.OMGD.GCP.Project != "" {
		if c.OMGD.GCP.CredsFile == "" {
			configDir := ""

			if runtime.GOOS == "windows" {
				configDir, err = os.UserConfigDir()

				if err != nil {
					LogFatal(fmt.Sprintf("Error finding user's config directory %s", err))
				}
			} else {
				homeDir, err := os.UserHomeDir()

				if err != nil {
					LogFatal(fmt.Sprintf("Error finding user's home directory %s", err))
				}

				configDir = fmt.Sprintf("%s/.config", homeDir)
			}

			c.OMGD.GCP.CredsFile = filepath.Join(configDir, "gcloud", "application_default_credentials.json")
		}
	}

	return &c
}

func GetProfileFromDir(path string, dir string) *ProfileConf {
	root, err := os.Getwd()
	if err != nil {
		LogFatal(fmt.Sprint(err))
	}

	err = os.Chdir(filepath.Join(root, dir))
	if err != nil {
		LogFatal(fmt.Sprint(err))
	}

	profile := GetProfile(path)
	profile.rootDir = dir

	err = os.Chdir(root)
	if err != nil {
		LogFatal(fmt.Sprint(err))
	}

	return profile
}

func (profile ProfileConf) SaveProfileFromMap(profileMap *map[interface{}]interface{}) {
	yamlBytes, err := yaml.Marshal(profileMap)

	if err != nil {
		LogFatal("Error marshalling from data to saving profile to yaml!")
	}

	profilePath := profile.path

	if profile.rootDir != "" {
		profilePath = filepath.Join(profile.rootDir, profilePath)
	}

	err = ioutil.WriteFile(profilePath, yamlBytes, 0755)

	if err != nil {
		LogFatal(fmt.Sprintf("Error on file write to saving profile to yaml! >> %s", err))
	}

	profile.LoadProfile()
}

func (profile ProfileConf) LoadProfile() *ProfileConf {
	var bytes, err = yaml.Marshal(profile.GetProfileAsMap())
	if err != nil {
		LogFatal(fmt.Sprintf("YAML marshal err from map: %v", err))
	}

	err = yaml.Unmarshal(bytes, &profile)
	if err != nil {
		LogFatal(fmt.Sprintf("YAML Unmarshal err: %v", err))
	}

	return &profile
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
					return
				} else {
					setValueToKeyWithArray(keys, keyIndex+1, v.(map[interface{}]interface{}), value)
					return
				}
			}
		}
	}

	// no match?
	if keyIndex == len(keys)-1 {
		obj[keys[keyIndex]] = value
	} else {
		obj[keys[keyIndex]] = map[interface{}]interface{}{}
		setValueToKeyWithArray(keys, keyIndex+1, obj[keys[keyIndex]].(map[interface{}]interface{}), value)
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
					LogFatal(fmt.Sprint("Key [", k, "] is map/list of values in one yaml and a singular value in the other yaml, can't merge them"))
				}
			} else {
				// key is not a map, so we just need to copy the value if they are both non-map types
				if !isMap(reflect.TypeOf((*merge)[k]).String()) {
					(*merge)[k] = v
				} else {
					LogFatal(fmt.Sprint("Key [", k, "] is map/list of values in one yaml and a singular value in the other yaml, can't merge them"))
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
