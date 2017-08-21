package uname

import (
	"bytes"
	"os/exec"
	"strings"
)

//Add Functionality that standard Uname syscall does not have.

const uname_bin string = "/usr/bin/uname"

func GetHardwarePlatform() string{
	return execUname("-i")
}

func GetProcessorType() string {
	return execUname("-p")
}

func execUname(arg string) string {
	uname := exec.Command(uname_bin, arg)
	var out bytes.Buffer
	uname.Stdout = &out
	err := uname.Run()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(out.String())
}