package cmd

import (
	"os"

	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/toasterson/mozaik/util"
	"github.com/toasterson/opencloud/image"
)

var (
	imgtype string
	imgdir string
	sets []string
	notresolve bool
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an image with the following options",
	Long: `Instructs imageadm to create an image as specified.
	Use -t to specify a type of image to create Supported Types are Chroot|ZFS|UFS
	`,
	Run: createCmdrun,
	Args: cobra.MinimumNArgs(1),
}

func init() {
	RootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&imgtype,"type", "t", "chroot", "Set this to Chroot, ZFS or UFS to specify which image to be created.")
	createCmd.Flags().StringVarP(&imgdir,"dir", "D", ".", "Which directory to create the image in defaults to ")
	createCmd.Flags().StringArrayVarP(&sets, "set", "s", []string{}, "Which Filesets to include in the final Image. use config to show the sets")
	createCmd.Flags().BoolVarP(&notresolve, "resolve", "n", false, "If the Files should be resolved from the sets. Defaults to true")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func createCmdrun(cmd *cobra.Command, args []string) {
	imgname := args[0]
	imgdir, err := filepath.Abs(imgdir)
	util.Must(err)
	imgp := filepath.Join(imgdir, imgname)
	err = os.Mkdir(imgp, 0755)
	if !os.IsExist(err){
		util.Must(err)
	}
	profile, err := image.NewProfile(imgname)
	util.Must(err)
	config, err := image.LoadConfiguration(cfgFile)
	for _, set := range sets{
		section, ok := config.Sections[set]
		if ok{
			profile.FileSets = append(profile.FileSets, set)
			profile.Users = append(profile.Users, config.GetAllFromSection(&section, "users")...)
			profile.Groups = append(profile.Groups, config.GetAllFromSection(&section, "groups")...)
			profile.Devices = append(profile.Devices, config.GetAllFromSection(&section, "devices")...)
		}
	}
	if notresolve == false {
		profile.ResolveFiles(&config)
	}
	err = profile.Save(imgp)
	util.Must(err)
}
