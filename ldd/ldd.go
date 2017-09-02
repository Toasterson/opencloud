package ldd

import (
	"debug/elf"

	"path/filepath"

	"io/ioutil"

	"strings"

	"os"

	"github.com/h2non/filetype"
	_ "github.com/h2non/filetype"
	"github.com/toasterson/opencloud/common"
)

var lib_paths = []string{"/usr/lib", "/lib"}
var bin_paths = []string{"/bin", "/usr/bin", "/sbin", "/usr/sbin"}
var arch_parts = []string{"amd64", "i86"}

func GetSharedLibraries(binary string, searchLocation []string) []string {
	lib_paths = append(lib_paths, searchLocation...)
	file, err := elf.Open(binary)
	if err != nil {
		panic(err)
	}
	libs := recListImportedLibs(file)
	common.RemoveDuplicates(&libs)
	common.RemoveEmpties(&libs)
	return libs
}

func FindLibrary(lib string) string{
	if !strings.Contains(lib, ".so"){
		lib = lib + ".so"
	}
	for _, libp := range lib_paths{
		for _, arch_part := range arch_parts{
			path := filepath.Join(libp, arch_part, lib)
			if common.FileExists(path){
				return path
			}
		}
		path := filepath.Join(libp, lib)
		if common.FileExists(path){
			return path
		}
	}
	return ""
}

func FindLibraries(lib string) []string{
	for _, libp := range lib_paths {
		path := filepath.Join(libp, lib)
		found, err := filepath.Glob(path)
		if err == nil && found != nil {
			return found
		}
	}
	return nil
}

func IsExecutableBinary(bin string) bool {
	buf, err := ioutil.ReadFile(bin)
	if err != nil {
		return false
	}
	return filetype.IsMIME(buf, "application/x-executable")
}

func FindBinary(bin string) string{
	paths := bin_paths
	path_var := os.ExpandEnv("$PATH")
	if strings.Contains(path_var, ":"){
		paths = strings.SplitN(path_var, ":", -1)
	}
	for _, binp := range paths{
		for _, arch_part := range arch_parts{
			path := filepath.Join(binp, arch_part, bin)
			if common.FileExists(path){
				return path
			}
		}
		path := filepath.Join(binp, bin)
		if common.FileExists(path){
			return path
		}
	}
	return ""
}

func recListImportedLibs(file *elf.File) []string {
	libs := []string{}
	syms, err := file.ImportedLibraries()
	if err != nil {
		panic(err)
	}
	for _, sym := range syms{
		libFilePath := FindLibrary(sym)
		libFile, err := elf.Open(libFilePath)
		if err != nil {
			continue
		}
		libs = append(libs, libFilePath)
		libs = append(libs, recListImportedLibs(libFile)...)
	}
	return libs
}


