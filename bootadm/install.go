package bootadm

import (
	"os/exec"
)

const bootadm_bin = "/sbin/bootadm"

func InstallBootLoader(rootDir string, pool string) error {
	args := []string{}
	if rootDir != "" {
		args = append(args, "-R", rootDir)
	}
	if pool != "" {
		args = append(args, "-P", pool)
	}
	return execBootadmInstall(args)
}

func execBootadmInstall(args []string) error {
	realArgs := []string{"install-bootloader"}
	realArgs = append(realArgs, args...)
	bootadm := exec.Command(bootadm_bin, realArgs...)
	return bootadm.Run()
}