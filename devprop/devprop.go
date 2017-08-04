package devprop

import (
	"os/exec"
	"bytes"
)

const devprop_bin string = "/sbin/devprop"

func GetValue(key string)  (value string, err error){
	cmd := exec.Command(devprop_bin, key)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err = cmd.Run(); err != nil{
		return
	}
	value = string(out.Bytes())
	return
}