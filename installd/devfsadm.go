package installd

import "os/exec"

var devfsadm_bin string = "/usr/sbin/devfsadm"

func CreateDeviceLinks(root string, device_classes []string) (err error){
	args := []string{}
	if root != ""{
		args = append(args, "-r", root)
	}
	for _, dev_class := range device_classes{
		args = append(args, "-c", dev_class)
	}
	devfsadm := exec.Command(devfsadm_bin, args...)
	if err = devfsadm.Run(); err != nil {
		return
	}
	return
}
