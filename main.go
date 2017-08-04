package main

import (
	"flag"
	"github.com/toasterson/opencloud/devprop"
	"strings"
	"github.com/toasterson/opencloud/installd"
	"os"
	"encoding/json"
	"github.com/toasterson/mozaik/util"
	"github.com/toasterson/mozaik/logger"
	"github.com/toasterson/opencloud/common"
)

func main() {
	confFile := flag.String("file", "NONE", "The Location of the Install Configuration file.")
	dryRun := flag.Bool("n", false, "Only report what would be done.")
	flag.Parse()
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
	conf, err := os.OpenFile(*confFile, os.O_RDONLY, 0444)
	defer conf.Close()
	util.Must(err)
	var buffer []byte
	conf.Read(buffer)
	var confObj installd.InstallConfiguration
	util.Must(json.Unmarshal(buffer, &confObj))
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