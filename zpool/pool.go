package zpool

import (
	"regexp"
	"strings"

	"github.com/toasterson/opencloud/zfs"
)

func OpenPool(name string) (p Pool) {
	//TODO switch to use -Hp as this does not print first line
	retVal, err := zpoolExec([]string{"get", "all", name})
	if err != nil {
		return
	}
	p.Name = name
	p.Properties = make(zfs.DatasetProperties)
	for _, line := range retVal {
		propLine := strings.Fields(line)
		propName := propLine[1]

		prop := zfs.Property{
			propLine[2],
			propLine[3],
		}
		p.Properties[propName] = prop

	}
	children, err := zfs.List(name)
	if err != nil {
		return
	}
	for _, child := range children {
		if !(child == name) {
			slash := regexp.MustCompile("/")
			matches := slash.FindAllStringIndex(child, -1)
			//zfs command outputs all Children But that is a hassle to parse so ignore children of children here
			//TODO Figure out if I want to switch this to nonrecursive. and if So How
			if !(len(matches) > 1 ) {
				p.Datasets = append(p.Datasets, zfs.OpenDataset(child))
			}
		}
	}
	return
}
