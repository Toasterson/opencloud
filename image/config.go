package image

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"path/filepath"

	"os"

	"github.com/toasterson/opencloud/common"
	"github.com/toasterson/opencloud/ldd"
)

var Default_path string = "/etc/imagedefs.json"
var tmp_transportvar_walkDir = []string{}
//var Default_path string = "$HOME/.config/imagedefs.json"

func init() {
	if strings.Contains(Default_path, "$"){
		Default_path = os.ExpandEnv(Default_path)
	}
}

func walkIntoTmpVar(path string, info os.FileInfo, err error) error{
	if !info.IsDir(){
		tmp_transportvar_walkDir = append(tmp_transportvar_walkDir, path)
	}
	return nil
}

type Config struct {
	Sections map[string]ConfigSection `json:"sections"`
}

type ConfigSection struct {
	Name string `json:"name"`
	Devices []string `json:"devices,omitempty"`
	Users []string `json:"users,omitempty"`
	Groups []string `json:"groups,omitempty"`
	Comment string `json:"comment,omitempty"`
	Paths []string `json:"paths,omitempty"`
	Dependencies []string `json:"dependencies,omitempty"`
}

func LoadConfiguration(path string) (conf Config, err error){
	if path == "" {
		path = Default_path
	}
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(buf, &conf)
	return
}

func SaveConfigSection(section ConfigSection) error {
	conf, err := LoadConfiguration("")
	if err != nil {
		conf = Config{}
		conf.Sections = make(map[string]ConfigSection)
	}
	conf.Sections[section.Name] = section
	confJson, err := json.Marshal(conf)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(Default_path, confJson, 0644)
}

func (c Config)GetFiles(sections []string) []string{
	files := []string{}
	for _, section := range sections{
		sectionObj := c.Sections[section]
		paths := c.GetAllFromSection(&sectionObj, "paths")
		for _, path := range paths{
			if strings.Contains(path, "*") {
				if !strings.Contains(path, "/"){
					libs := ldd.FindLibraries(path)
					if libs != nil {
						files = append(files, libs...)
					}
				} else {
					found, err := filepath.Glob(path)
					if err == nil && found != nil {
						for _, foundp := range found{
							foundStat, err := os.Stat(foundp)
							if err != nil {
								continue
							}
							if foundStat.Mode().IsDir(){
								tmp_transportvar_walkDir = []string{}
								filepath.Walk(foundp, walkIntoTmpVar)
								files = append(files, tmp_transportvar_walkDir...)
								tmp_transportvar_walkDir = []string{}
								continue
							}
							files = append(files, foundp)
						}
					}
				}
			} else if !strings.Contains(path, "/"){
				if strings.Contains(path, "lib"){
					files = append(files, ldd.FindLibrary(path))
				} else {
					files = append(files, ldd.FindBinary(path))
				}
			} else {
				pStat, err := os.Stat(path)
				if err != nil {
					continue
				}
				if pStat.Mode().IsDir(){
					tmp_transportvar_walkDir = []string{}
					filepath.Walk(path, walkIntoTmpVar)
					files = append(files, tmp_transportvar_walkDir...)
					tmp_transportvar_walkDir = []string{}
					continue
				}
				files = append(files, path)
			}
		}
		for _, file := range files {
			if ldd.IsExecutableBinary(file){
				files = append(files, ldd.GetSharedLibraries(file, []string{})...)
			}
		}
	}
	common.RemoveDuplicates(&files)
	common.RemoveEmpties(&files)
	return files
}

func (c Config) GetAllFromSection(section *ConfigSection, variable string) []string {
	var retVal []string
	switch variable {
	case "users":
		retVal = section.Users
	case "groups":
		retVal = section.Groups
	case "devices":
		retVal = section.Devices
	case "paths":
		retVal = section.Paths
	}
	for _, dep := range section.Dependencies{
		subsec, ok := c.Sections[dep]
		if ok {
			retVal = append(retVal, c.GetAllFromSection(&subsec, variable)...)
		}
	}
	return retVal
}