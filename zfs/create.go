package zfs

// CreateDataset create a new filesystem or volume on path representing
// pool/dataset
func CreateDataset(path string, dtype DatasetType, props map[string]string, createOnly bool) (d Dataset, err error) {
	args := []string{"create"}
	if dtype == DatasetTypeVolume {
		dsize := props["size"]
		dblocksize := props["blocksize"]
		delete(props, "blocksize")
		delete(props, "size")
		args = append(args, "-V", dsize)
		if dblocksize != "" {
			args = append(args, "-b", dblocksize)
		}
	}
	args = append(args, DatasetPropertyListToCMD(props)...)
	args = append(args, path)
	if _, err = zfsExec(args); err != nil {
		return
	}
	if !createOnly {
		d = OpenDataset(path)
	}
	return
}
