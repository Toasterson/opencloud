package zfs

import (
	"os/exec"
	"bytes"
	"strings"
	"github.com/c2h5oh/datasize"
	"errors"
)

func List(zpool string) (datasets []string, err error) {
	datasets, err = zfsList([]string{"-r", "-o", "name", zpool})
	return datasets, err
}

func Size(dataset string) (size datasize.ByteSize, err error) {
	return zfsListSomeSize(dataset, "referenced")
}

func Avail(dataset string) (size datasize.ByteSize, err error) {
	return zfsListSomeSize(dataset, "available")
}

func Used(dataset string) (size datasize.ByteSize, err error) {
	return zfsListSomeSize(dataset, "usedbydataset")
}

func UsedIncludingChildren(dataset string) (size datasize.ByteSize, err error) {
	return zfsListSomeSize(dataset, "used")
}

func zfsListSomeSize(dataset string, parameters ...string) (size datasize.ByteSize, err error) {
	//TODO switch to use -Hp as this does not print first line
	if dataset == "" {
		return size, errors.New("Dataset is not allowed to be empty.")
	}
	zfs_args := []string{"-o"}
	for i, param := range parameters {
		if i >= 1 {
			zfs_args = append(zfs_args, ","+param)
		} else {
			zfs_args = append(zfs_args, param)
		}
	}
	zfs_args = append(zfs_args, dataset)
	datasetSize, err := zfsList(zfs_args)
	if err != nil {
		return
	}
	return convertToSize(datasetSize[0])
}

func zfsList(args []string) (retVal []string, err error) {
	args = append([]string{"list"}, args...)
	cmd := exec.Command("zfs", args...)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err = cmd.Run(); err != nil {
		return retVal, err
	}
	retVal = strings.Split(out.String(), "\n")
	retVal = retVal[1:]
	//Do some trimming as there could be a empty line in there
	for i, val := range retVal {
		val = strings.TrimSpace(val)
		if val == "" {
			retVal = append(retVal[:i], retVal[i+1:]...)
		}
	}
	return retVal, nil
}
