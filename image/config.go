package image

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"os"

	"github.com/toasterson/glog"
	"github.com/toasterson/opencloud/common"
	"github.com/toasterson/opencloud/ldd"
	"github.com/toasterson/uxfiletool"
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
		glog.Debugf("Getting Files for %s", sectionObj.Name)
		paths := c.GetAllFromSection(&sectionObj, "paths")
		glog.Tracef("Found %v", paths)
		for _, path := range paths{
			if strings.Contains(path, "*") {
				//Lets assume we have a Pattern to resolve
				glog.Debugf("%s needs globbing", path)
				found := uxfiletool.FindByGlob(path)
				glog.Tracef("Found %v", found)
				files = append(files, found...)
			} else if strings.Contains(path, "/"){
				//A / in the path means the path needs to be full, thus add it if its a file and add dir contents
				pStat, err := os.Lstat(path)
				if err != nil {
					glog.Warnf("%s can not be read ignoring", path)
					continue
				}
				if pStat.Mode().IsDir(){
					glog.Debugf("%s is a directory getting Files", path)
					foundfiles := uxfiletool.FindAllIn(path, uxfiletool.FindTypeFile)
					glog.Tracef("Found files: %v", foundfiles)
					foundlinks := uxfiletool.FindAllIn(path, uxfiletool.FindTypeLink)
					glog.Tracef("Found links: %v", foundlinks)
					files = append(files, foundfiles...)
					files = append(files, foundlinks...)
				} else {
					glog.Debugf("%s is a File", path)
					files = append(files, path)
				}
			} else {
				//Lastly assume be will find that binary/lib in PATH
				glog.Debugf("Assuming you want me to resolve %s in Path", path)
				if strings.Contains(path, "lib"){
					libs := uxfiletool.FindLib(path)
					glog.Tracef("Found: %v libs", libs)
					files = append(files, libs...)
				} else {
					bins := uxfiletool.FindInPath(path)
					glog.Tracef("Found the following Binaries %v", bins)
					files = append(files, bins...)
				}
			}
		}
	}
	common.RemoveDuplicates(&files)
	common.RemoveEmpties(&files)
	//After we have resolved all the relativity lets grab the shared libs of the files
	for _, file := range files {
		if ldd.IsExecutableBinary(file){
			glog.Debugf("Getting shared libs of: %s", file)
			libs := ldd.GetSharedLibraries(file, []string{})
			glog.Tracef("Found: %v", libs)
			files = append(files, libs...)
		}
	}
	//After all the filling we will have duplicates and empty entries
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