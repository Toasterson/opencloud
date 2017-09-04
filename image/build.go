package image

import (
	"path/filepath"

	"archive/tar"
	"compress/gzip"
	"os"

	"github.com/appc/spec/aci"
	"github.com/pkg/errors"
	"github.com/toasterson/glog"
	"github.com/toasterson/uxfiletool"
)

func BuildChroot(imageProfile *Profile, target string) error {
	for _, file := range imageProfile.Files {
		glog.Debugf("Copying: %s", file)
		if err := uxfiletool.ExactCopy(file, target); err != nil {
			glog.Tracef("Error Encountered: %s", err)
			if !os.IsNotExist(err){
				glog.Errf("Error: %s", err)
				return err
			}
		}
	}
	return nil
}

func BuildACI(imageProfile *Profile, target string) error {
	outFilePath := target
	if target == ""{
		outFilePath = filepath.Join(target, string(imageProfile.Manifest.Name), ".tar.gz")
	}
	outFile, err := os.Create(outFilePath)
	if err != nil{
		return err
	}
	gw := gzip.NewWriter(outFile)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	aciW := aci.NewImageWriter(imageProfile.Manifest, tw)
	defer aciW.Close()
	for _, file := range imageProfile.Files{
		fileObj, err := os.Open(file)
		fileStat, _ := fileObj.Stat()
		if err != nil{
			if os.IsNotExist(err){
				continue
			}
			return nil
		}
		var linkTarget string = ""
		if lt, err := os.Readlink(file); err != nil {
			linkTarget = lt
		}
		th, err := tar.FileInfoHeader(fileStat, linkTarget)
		if err != nil{
			return err
		}
		th.Name = filepath.Join("rootfs", file)
		if err := aciW.AddFile(th, fileObj); err != nil {
			return err
		}
	}
	return nil
}

func BuildUFS(imageProfile *Profile, target string) error {
	return errors.New("UFS Image Not Implemented")
}

func BuildZFS(imageProfile *Profile, target string) error {
	return errors.New("ZFS Image Not Implemented")
}

func BuildTar(imageProfile *Profile, target string) error {
	return errors.New("Tar Image Not Implemented")
}