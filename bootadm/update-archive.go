package bootadm

func UpdateBootArchive(rootDir string) error {
	args := []string{"update-archive"}
	if rootDir != "" && rootDir != "/"{
		args = append(args, "-R", rootDir)
	}
	return execBootadm(args)
}