package common

import (
	"errors"
	"fmt"
	"log"
	"os"
)

var Stdlog, Errlog *log.Logger

func NotSupportedError(functionality string) error {
	return errors.New(fmt.Sprintf("Functionality %s is currently not Supported", functionality))
}

func InvalidConfiguration(section string) error {
	return errors.New(fmt.Sprintf("Invalid Configuration %s is not correct please fix", section))
}

func init() {
	Stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	Errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}
