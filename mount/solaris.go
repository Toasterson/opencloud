package mount

import (
	"os/exec"
	"bytes"
	"fmt"
	"strings"
	"errors"
)

const (
	lofiadm_bin string = "/usr/sbin/lofiadm"
	mount_bin string = "/usr/sbin/mount"
)

func MountLoopDevice(fstype string, mountpoint string, file string) error {
	lofiadm := exec.Command(lofiadm_bin, "-a", file)
	var lofibuff, lofierr bytes.Buffer
	lofiadm.Stdout = lofibuff
	lofiadm.Stderr = lofierr
	if err := lofiadm.Run(); err != nil{
		return errors.New(strings.TrimSpace(lofierr.String()))
	}
	var mountbuff, mounterr bytes.Buffer
	mount := exec.Command(mount_bin, fmt.Sprintf("-F%s", fstype), "-o", "ro", lofibuff.String(), mountpoint)
	mount.Stdout = mountbuff
	mount.Stderr = mounterr
	if err := mount.Run(); err != nil{
		return errors.New(strings.TrimSpace(mounterr.String()))
	}
	return nil
}
