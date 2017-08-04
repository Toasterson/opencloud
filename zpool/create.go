package zpool

import "github.com/toasterson/opencloud/zfs"

// Create a new Zfs Pool on the Followin disks with mode
func CreatePool(name string, props map[string]string, force bool, poolType string, drives []string, createOnly bool) (p Pool, err error) {
	args := []string{"create"}
	if force{
		args = append(args, "-f")
	}
	args = append(args, zfs.DatasetPropertyListToCMD(props)...)
	args = append(args, name)
	if poolType != ""{
		args = append(args, poolType)
	}
	args = append(args, drives...)
	if _, err = zpoolExec(args); err != nil{
		return
	}
	if !createOnly {
		p = OpenPool(name)
	}
	return
}
