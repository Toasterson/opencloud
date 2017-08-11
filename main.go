package main

import (
	"flag"
	"github.com/toasterson/opencloud/devprop"
	"github.com/toasterson/opencloud/installd"
	"encoding/json"
	"github.com/toasterson/mozaik/util"
	"github.com/toasterson/mozaik/logger"
	"io/ioutil"
)

func main() {
	confFile := flag.String("file", "NONE", "The Location of the Install Configuration file.")
	dryRun := flag.Bool("n", false, "Only report what would be done.")
	flag.Parse()
	/*
	//First Try to load the config file via http
	configLocation , err := devprop.GetValue("install_config")
	//Assume we have gotten a property if err is nil TODO Make sure we get error if we do not have a value
	if err == nil {
		if strings.HasPrefix(configLocation, "http") || strings.HasPrefix(configLocation, "https") {
			var err error
			*confFile, err = installd.HTTPDownload(configLocation, "")
			util.Must(err)
		} else if strings.HasPrefix(configLocation, "nfs"){
			panic(common.NotSupportedError("Nfs Downloads"))
		}
	}
	*/
	file, err := ioutil.ReadFile(*confFile)
	util.Must(err)
	var confObj installd.InstallConfiguration
	util.Must(json.Unmarshal(file, &confObj))
	if confObj.MediaURL == ""{
		confObj.MediaURL, err = devprop.GetValue("install_media")
		util.Must(err)
	}
	if !*dryRun {
		installd.Install(confObj)
	} else {
		//TODO more detailed Dry-Run reporting
		logger.Info(confObj)
	}
}