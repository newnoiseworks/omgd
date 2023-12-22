package utils

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func getData(profile *ProfileConf, buildPath string) *map[interface{}]interface{} {
	data := make(map[interface{}]interface{})
	data["profile"] = profile.GetProfileAsMap()

	resourceDir := buildPath

	if strings.HasPrefix(profile.path, "..") {
		paths := strings.Split(profile.path, fmt.Sprintf(".%s", string(os.PathSeparator)))

		resourceDir = ""

		for _, p := range paths {
			if p == "." {
				resourceDir += fmt.Sprintf("..%s", string(os.PathSeparator))
			}
		}
	}

	resourceDir = filepath.Join(resourceDir, "resources")
	_, err := os.Stat(resourceDir)

	if !os.IsNotExist(err) {
		err = filepath.Walk(resourceDir, func(tmpl string, info fs.FileInfo, err error) error {
			if err != nil {
				LogDebug(fmt.Sprintf("no resources directory found in %v", resourceDir))
				return nil
			}

			name := info.Name()

			if info.IsDir() == false && strings.HasSuffix(name, ".yml") {
				c := make(map[interface{}]interface{})

				yamlFile, err := ioutil.ReadFile(tmpl)
				if err != nil {
					LogDebug(fmt.Sprintf("yamlFile Get err: #%v ", err))
				}

				err = yaml.Unmarshal(yamlFile, &c)
				if err != nil {
					LogDebug(fmt.Sprintf("Unmarshal err: %v", err))
				}

				for k, v := range c {
					data[k] = v
				}
			}

			return nil
		})

		if err != nil {
			LogDebug(fmt.Sprint(err))
		}
	}

	return &data
}

func BuildTemplatesFromPath(profile *ProfileConf, buildPath string, templateExtension string, removeTemplateAfterProcessing bool) {
	data := getData(profile, buildPath)

	LogDebug(fmt.Sprintf("building template files in %s", buildPath))

	err := filepath.Walk(buildPath, func(tmpl string, info fs.FileInfo, err error) error {
		if err != nil {
			LogFatal(fmt.Sprint(err))
		}

		name := info.Name()

		if info.IsDir() == false && strings.HasSuffix(name, "."+templateExtension) {
			processTemplate(tmpl, data, templateExtension, removeTemplateAfterProcessing)
		}

		return nil
	})

	if err != nil {
		LogFatal(fmt.Sprint(err))
	}
}

func BuildTemplateFromPath(tmplPath string, profile *ProfileConf, buildPath string, templateExtension string, removeTemplateAfterProcessing bool) {
	data := getData(profile, buildPath)
	processTemplate(tmplPath, data, templateExtension, removeTemplateAfterProcessing)
}

func processTemplate(tmpl string, data *map[interface{}]interface{}, templateExtension string, removeTemplateAfterProcessing bool) {
	final_path := strings.ReplaceAll(tmpl, "."+templateExtension, "")

	LogDebug(fmt.Sprintf("processing template file %s >> %s", tmpl, final_path))

	t := template.New(path.Base(tmpl)).Funcs(template.FuncMap{
		"md5":             StrToMd5,
		"upperSnake":      StrToUpperSnake,
		"camel":           StrToCamel,
		"gcpZoneToRegion": GCPZoneToRegion,
	})

	if templateExtension == "omgdtpl" {
		t.Delims("{*", "*}")
	}

	t, err := t.ParseFiles(tmpl)

	if err != nil {
		LogFatal(fmt.Sprint(err))
	}

	f, err := os.Create(final_path)
	if err != nil {
		LogFatal(fmt.Sprint(err))
	}

	err = t.Execute(f, data)
	if err != nil {
		LogFatal(fmt.Sprint(err))
	}

	if removeTemplateAfterProcessing {
		err = os.Remove(tmpl)
		if err != nil {
			LogFatal(fmt.Sprint(err))
		}
	}
}
