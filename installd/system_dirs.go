// +build solaris,cgo

package installd

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"

	"github.com/toasterson/mozaik/logger"
)

type DirConfig struct {
	Name string
	Mode int
	Owner string
	Group string
}

var defualtDirectories = []DirConfig{
	{Name: "tmp", Mode: 1777},
	{Name: "system", Mode: 555},
	{Name: "system/contract", Mode: 555},
	{Name: "system/object", Mode: 555},
	{Name: "system/boot", Mode: 555},
	{Name: "proc", Mode: 555},
	{Name: "opt", Group: "sys"},
	{Name: "mnt", Group: "sys"},
	{Name: "devices/", Group: "sys"},
	{Name: "devices/pseudo", Group: "sys"},
	{Name: "dev", Group: "sys"},
	{Name: "dev/fd", Group: "sys"},
	{Name: "dev/rmt", Group: "sys"},
	{Name: "dev/swap", Group: "sys"},
	{Name: "dev/dsk", Group: "sys"},
	{Name: "dev/rdsk", Group: "sys"},
	{Name: "dev/net", Group: "sys"},
	{Name: "dev/ipnet", Group: "sys"},
	{Name: "dev/sad", Group: "sys"},
	{Name: "dev/pts", Group: "sys"},
	{Name: "dev/term", Group: "sys"},
	{Name: "dev/vt", Group: "sys"},
	{Name: "dev/zcons", Group: "sys"},
}

func MakeSystemDirectories(rootDir string, dirs []DirConfig){
	dirs = append(dirs, defualtDirectories...)
	for _, dir := range dirs {
		path := fmt.Sprintf("%s/%s", rootDir, dir.Name)
		logger.Trace(fmt.Sprintf("Creating System Directory %s", path))
		var uid, gid int
		os.Mkdir(path, 0755)
		if dir.Mode != 0 {
			syscall.Chmod(path, uint32(dir.Mode))
		}
		if dir.Owner != "" {
			owner, err := user.Lookup(dir.Owner)
			if err != nil {
				logger.Error(fmt.Sprintf("User %s does not exist this should not happen %s", dir.Owner, err))
			} else {
				uid, _ = strconv.Atoi(owner.Uid)
			}
		}
		if dir.Group != "" {
			group, err := user.LookupGroup(dir.Group)
			if err != nil {
				logger.Error(fmt.Sprintf("Group %s does not exist this should not happen: %s", dir.Group, err))
			} else {
				gid, _ = strconv.Atoi(group.Gid)
			}
		}
		dirMode, _ := os.Stat(path)
		dirStat := dirMode.Sys().(*syscall.Stat_t)
		if uid == 0 {
			uid = int(dirStat.Uid)
		}
		if gid == 0 {
			gid = int(dirStat.Gid)
		}
		syscall.Chown(path, uid, gid)
	}
}
