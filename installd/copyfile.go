package installd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"github.com/toasterson/mozaik/logger"
	"github.com/toasterson/mozaik/util"
)

func walkCopy(path string, info os.FileInfo, err error) error {
	dstpath := strings.Replace(path, altMountLocation, altRootLocation, -1)
	lsrcinfo, err := os.Lstat(path)
	if os.IsNotExist(err) {
		//Ignore Unexistent Directories
		return nil
	}
	util.Must(err)
	if info.IsDir() {
		logger.Trace(fmt.Sprintf("Mkdir %s", dstpath))
		util.Must(os.Mkdir(dstpath, info.Mode()))
	} else if lsrcinfo.Mode()&os.ModeSymlink != 0 {
		//We have a Symlink thus Create it on the Target
		dstTarget, _ := os.Readlink(path)
		logger.Trace(fmt.Sprintf("Creating Symlink %s -> %s", dstpath, dstTarget))
		util.Must(os.Symlink(dstTarget, dstpath))
	} else {
		//We Have a regular File Copy it
		go copyFileExact(path, info, dstpath)
	}
	return nil
}

func copyFileExact(source string, srcInfo os.FileInfo, dest string) {
	//logger.Trace(fmt.Sprintf("Copy %s -> %s", path, dest))
	src, err := os.Open(source)
	defer src.Close()
	util.Must(err)
	dst, err := os.Create(dest)
	defer dst.Close()
	util.Must(err)
	_, err = io.Copy(dst, src)
	util.Must(err)
	//util.Must(dst.Sync())
	srcStat := srcInfo.Sys().(*syscall.Stat_t)
	util.Must(syscall.Chmod(dest, srcStat.Mode))
	util.Must(syscall.Chown(dest, int(srcStat.Uid), int(srcStat.Gid)))
	//util.Must(os.Chtimes(dest, time.Unix(int64(srcStat.Atim.Sec),int64(srcStat.Atim.Nsec)), time.Unix(int64(srcStat.Mtim.Sec),int64(srcStat.Mtim.Nsec))))
}