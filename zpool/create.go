package zpool

import (
	"github.com/toasterson/opencloud/zfs"
	"runtime"
	"fmt"
)

// Create a new Zfs Pool on the Followin disks with mode
func CreatePool(name string, props map[string]string, force bool, poolType string, drives []string, createOnly bool) (p Pool, err error) {
	args := []string{"create"}
	if force{
		args = append(args, "-f")
	}
	args = append(args, zfs.DatasetPropertyListToCMD(props)...)
	args = append(args, name)
	if poolType != "normal" && poolType != ""{
		args = append(args, poolType)
	}
	for _,drive := range drives {
		args = append(args, fmt.Sprintf("%s/%s", drivePath(), drive))
	}
	if _, err = zpoolExec(args); err != nil{
		return
	}
	if !createOnly {
		p = OpenPool(name)
	}
	return
}

func drivePath() string{
	switch runtime.GOOS {
	case "linux":
		return "/dev"
	case "solaris":
		return "/dev/dsk"
	default:
		panic("Not Supported Os")
	}
	return ""
}