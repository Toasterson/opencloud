package zfs

import (
	"os/exec"
	"bytes"
	"strings"
	"github.com/c2h5oh/datasize"
	"strconv"
	"errors"
)

func List(zpool string) (datasets []string, err error) {
	datasets, err = zfsList([]string{"-r", "-o", "name", zpool})
	return datasets, err
}

func Size(dataset string) (size datasize.ByteSize, err error){
	return zfsListSomeSize(dataset, "referenced")
}

func Avail(dataset string) (size datasize.ByteSize, err error){
	return zfsListSomeSize(dataset, "available")
}

func Used(dataset string) (size datasize.ByteSize, err error){
	return zfsListSomeSize(dataset, "usedbydataset")
}

func UsedIncludingChildren(dataset string) (size datasize.ByteSize, err error){
	return zfsListSomeSize(dataset, "used")
}

func zfsListSomeSize(dataset string, parameters ...string) (size datasize.ByteSize, err error){
	if dataset == "" {
		return size, errors.New("Dataset is not allowed to be empty.")
	}
	zfs_args := []string{"-o"}
	for i, param := range(parameters){
		if i >= 1{
			zfs_args = append(zfs_args, ","+param)
		} else {
			zfs_args = append(zfs_args, param)
		}
	}
	zfs_args = append(zfs_args, dataset)
	datasetSize, err := zfsList(zfs_args)
	if err != nil {
		return size, err
	}
	sizeText := strings.TrimSpace(datasetSize[0])
	if strings.Contains(sizeText, "."){
		unit := sizeText[len(sizeText)-1:]
		sizeText = sizeText[0:len(sizeText)-1]
		switch unit {
		case "T":
			unit = "G"
		case "G":
			unit = "M"
		case "M":
			unit = "K"
		default:
			unit = ""
		}
		f, ferr := strconv.ParseFloat(sizeText, 64)
		if ferr != nil {
			return size, ferr
		}
		f = f * 1024
		sizeText = strconv.FormatFloat(f, 'f', 0, 64) + unit
	}
	if uerr := size.UnmarshalText([]byte(sizeText)); uerr != nil{
		return size, uerr
	}
	return size, nil
}

func zfsList(args []string) (retVal []string, err error){
	args = append([]string{"list"}, args...)
	cmd := exec.Command("zfs", args...)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err = cmd.Run(); err != nil{
		return retVal, err
	}
	retVal = strings.Split(string(out.Bytes()), "\n")
	retVal = retVal[1:]
	//Do some trimming as there could be a empty line in there
	for i, val := range(retVal){
		val = strings.TrimSpace(val)
		if val == ""{
			retVal = append(retVal[:i], retVal[i+1:]...)
		}
	}
	return retVal, nil
}
