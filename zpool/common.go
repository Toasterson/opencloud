package zpool

import (
	"os/exec"
	"bytes"
	"strings"
	"errors"
	"github.com/toasterson/mozaik/logger"
	"fmt"
	"github.com/toasterson/opencloud/zfs"
)

// Pool - ZFS dataset object
type Pool struct {
	Name 	   string
	Properties zfs.DatasetProperties
	Datasets   []zfs.Dataset
}

func zpoolExec(args []string) (retVal []string, err error){
	logger.Trace(fmt.Sprintf("zpool %s", args))
	cmd := exec.Command("zpool", args...)
	var out, serr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &serr

	if err = cmd.Run(); err != nil{
		return []string{}, errors.New(strings.TrimSpace(serr.String()))
	}
	if out.Len() > 0 {
		retVal = strings.Split(out.String(), "\n")
		retVal = retVal[1:]
		//Do some trimming as there could be a empty line in there
		for i, val := range retVal {
			val = strings.TrimSpace(val)
			if val == ""{
				retVal = append(retVal[:i], retVal[i+1:]...)
			}
		}
	}
	return
}