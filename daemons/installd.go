package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"

	"strings"

	"github.com/toasterson/mozaik/logger"
	"github.com/toasterson/mozaik/util"
	"github.com/toasterson/opencloud/common"
	"github.com/toasterson/opencloud/devprop"
	"github.com/toasterson/opencloud/installd"
)

func main() {
	confFile := flag.String("c", "NONE", "The Location of the Install Configuration file. Can be http.")
	dryRun := flag.Bool("n", false, "Only report what would be done.")
	flag.Parse()

	//First Try to load the config file via http
	//Use confFile if
	if *confFile == "NONE" {
		*confFile = devprop.GetValue("install_config")
	}
	if strings.HasPrefix(*confFile, "http") || strings.HasPrefix(*confFile, "https") {
		var err error
		*confFile, err = installd.HTTPDownloadTo(*confFile, "")
		util.Must(err)
	} else if strings.HasPrefix(*confFile, "nfs"){
		panic(common.NotSupportedError("Nfs Downloads"))
	}
	file, err := ioutil.ReadFile(*confFile)
	util.Must(err)
	var confObj installd.InstallConfiguration
	util.Must(json.Unmarshal(file, &confObj))
	if confObj.MediaURL == "" {
		confObj.MediaURL = devprop.GetValue("install_media")
	}
	if !*dryRun {
		installd.Install(confObj)
	} else {
		//TODO more detailed Dry-Run reporting
		logger.Info(confObj)
	}
}
