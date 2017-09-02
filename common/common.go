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

func RemoveDuplicates(xs *[]string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *xs {
		if !found[x] {
			found[x] = true
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]
}

// Exists reports whether the named file or directory exists.
func FileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func RemoveEmpties(xs *[]string){
	j := 0
	for i, x := range *xs {
		if x != ""{
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]
}