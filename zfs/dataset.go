package zfs

import (
	"strings"
	"fmt"
	"regexp"
	"github.com/c2h5oh/datasize"
)

const (
	msgDatasetIsNil = "Dataset handle not initialized or its closed"
)

// DatasetProperties type is map of dataset or volume properties prop -> value
type DatasetProperties map[string]Property

// DatasetType defines enum of dataset types
type DatasetType int32

const (
	// DatasetTypeFilesystem - file system dataset
	DatasetTypeFilesystem DatasetType = 1 << 0
	// DatasetTypeSnapshot - snapshot of dataset
	DatasetTypeSnapshot = 1 << 1
	// DatasetTypeVolume - volume (virtual block device) dataset
	DatasetTypeVolume = 1 << 2
	// DatasetTypePool - pool dataset
	DatasetTypePool = 1 << 3
	// DatasetTypeBookmark - bookmark dataset
	DatasetTypeBookmark = 1 << 4
)

// Dataset - ZFS dataset object
type Dataset struct {
	Path       string
	Type       DatasetType
	Properties DatasetProperties
	Children   []Dataset
}

//Read a Dataset and all its Properties from zfs Command
func OpenDataset(path string) (d Dataset) {
	//TODO switch to use -Hp as this does not print first line
	retVal, err := zfsExec([]string{"get", "all", path})
	if err != nil {
		return
	}
	d.Path = path
	d.Properties = make(DatasetProperties)
	for _, line := range retVal {
		propLine := strings.Fields(line)
		propName := propLine[1]
		if propName == "type" {
			switch propLine[2] {
			case "filesystem":
				d.Type = DatasetTypeFilesystem
			case "volume":
				d.Type = DatasetTypeVolume
			default:
				d.Type = DatasetTypeFilesystem
			}
		} else {
			prop := Property{
				propLine[2],
				propLine[3],
			}
			d.Properties[propName] = prop
		}
	}
	children, err := List(path)
	if err != nil {
		return
	}
	for _, child := range children {
		if !(child == path) {
			slash := regexp.MustCompile("/")
			matches := slash.FindAllStringIndex(child, -1)
			//zfs command outputs all Children But that is a hassle to parse so ignore children of children here
			//TODO Figure out if I want to switch this to nonrecursive. and if So How
			if !(len(matches) > 1 ) {
				d.Children = append(d.Children, OpenDataset(child))
			}
		}
	}
	return
}

// SetProperty set ZFS dataset property to value. Not all properties can be set,
// some can be set only at creation time and some are read only.
// Always check if returned error and its description.
func (d *Dataset) SetProperty(prop string, value string) (err error) {
	if _, err = zfsExec([]string{"set", fmt.Sprintf("%s=%s", prop, value), d.Path}); err != nil {
		return
	}
	d.Properties[prop], err = d.GetProperty(prop)
	return
}

// GetProperty reload and return single specified property. This also reloads requested
// property in Properties map.
func (d *Dataset) GetProperty(p string) (prop Property, err error) {
	var retVal []string
	if retVal, err = zfsExec([]string{"get", p, d.Path}); err != nil {
		return
	}
	//TODO switch to use -Hp as this does not print first line
	propLine := strings.Fields(retVal[0])
	prop.Value = propLine[1]
	prop.Source = propLine[2]
	return
}

// Rename dataset
func (d *Dataset) Rename(newName string, forceUnmount bool) (err error) {
	args := []string{"rename"}
	if forceUnmount {
		args = append(args, "-f")
	}
	args = append(args, d.Path, newName)
	_, err = zfsExec(args)
	return
}

// IsMounted checks to see if the mount is active.  If the filesystem is mounted,
// sets in 'where' argument the current mountpoint, and returns true.  Otherwise,
// returns false.
func (d *Dataset) IsMounted() (mounted bool, where string) {
	if d.Properties["mounted"].Value == "yes" {
		mounted = true
		where = d.Properties["mountpoint"].Value
	} else {
		mounted = false
	}
	return
}

// Mount the given filesystem.
func (d *Dataset) Mount(options string) (err error) {
	args := []string{"mount"}
	if options != "" {
		args = append(args, "-o", options)
	}
	args = append(args, d.Path)
	_, err = zfsExec(args)
	return
}

// Unmount the given filesystem.
func (d *Dataset) Unmount() (err error) {
	_, err = zfsExec([]string{"unmount", d.Path})
	return
}

// UnmountAll unmount this filesystem and any children inheriting the
// mountpoint property.
func (d *Dataset) UnmountAll() (err error) {
	for _, child := range d.Children {
		if err = child.UnmountAll(); err != nil {
			return
		}
		if strings.Contains(child.Properties["mountpoint"].Source, "inherited") {
			if err = child.Unmount(); err != nil {
				return
			}
		}
	}
	return
}

func (d *Dataset) Size() (size datasize.ByteSize) {
	var err error
	if size, err = convertToSize(d.Properties["referenced"].Value); err != nil {
		return 0
	}
	return
}

func (d *Dataset) Avail() (size datasize.ByteSize) {
	var err error
	if size, err = convertToSize(d.Properties["available"].Value); err != nil {
		return 0
	}
	return
}

func (d *Dataset) Used() (size datasize.ByteSize) {
	var err error
	if size, err = convertToSize(d.Properties["usedbydataset"].Value); err != nil {
		return 0
	}
	return
}
