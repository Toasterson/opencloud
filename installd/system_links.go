package installd

import (
	"os"
	"fmt"
	"github.com/toasterson/mozaik/logger"
)

type LinkConfig struct {
	Name string
	Target string
}

var defaultLinks = []LinkConfig{
	{Name: "stderr", Target: "../fd/2"},
	{Name: "stdout", Target: "../fd/1"},
	{Name: "stdin", Target: "../fd/0"},
	{Name: "dld", Target: "../devices/pseudo/dld@0:ctl"},
}

func MakeDeviceLinks(rootDir string, links []LinkConfig){
	links = append(links, defaultLinks...)
	for _, link := range links {
		path := fmt.Sprintf("%s/dev/%s", rootDir, link.Name)
		err := os.Symlink(link.Target, path)
		if err != nil{
			logger.Error(err)
		}
	}
}
