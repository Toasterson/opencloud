package mount

import (
	"os/exec"
	"bytes"
	"fmt"
	"strings"
	"errors"
	"github.com/toasterson/mozaik/logger"
)

const (
	lofiadm_bin string = "/usr/sbin/lofiadm"
	mount_bin string = "/usr/sbin/mount"
)

func MountLoopDevice(fstype string, mountpoint string, file string) error {
	lofibuff, lofierr, err := lofiExec([]string{"-a", file})
	if err != nil{
		if !strings.Contains(lofierr, "Device busy") {
			return errors.New(strings.TrimSpace(lofierr))
		}
		//If we get Device Busy then we are already available in lofi. We do not need to fail here it is enough to fail
		//while mounting
		//Still we need to find out the device name
		lofibuff, lofierr, err = lofiExec([]string{file})
	}
	var mountbuff, mounterr bytes.Buffer
	mount := exec.Command(mount_bin, fmt.Sprintf("-F%s", fstype), "-o", "ro", lofibuff, mountpoint)
	logger.Trace(mount.Path, mount.Args)
	mount.Stdout = &mountbuff
	mount.Stderr = &mounterr
	if err := mount.Run(); err != nil{
		return errors.New(strings.TrimSpace(mounterr.String()))
	}
	return nil
}

//TODO func IsMounted(device, path) bool

func lofiExec(args []string) (out string, errout string, err error){
	lofiadm := exec.Command(lofiadm_bin, args...)
	var lofibuff, lofierr bytes.Buffer
	lofiadm.Stdout = &lofibuff
	lofiadm.Stderr = &lofierr
	logger.Trace(lofiadm.Path, lofiadm.Args)
	err = lofiadm.Run()
	out = strings.TrimSpace(lofibuff.String())
	errout = strings.TrimSpace(lofierr.String())
	return
}
