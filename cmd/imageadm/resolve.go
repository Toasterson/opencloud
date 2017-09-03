package cmd

import (
	"fmt"

	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/toasterson/opencloud/image"
)

var (
	resolve_save    bool
)

// buildCmd represents the create command
var resolveCmd = &cobra.Command{
	Use:   "resolve",
	Short: "Resolve the paths in profile",
	Long: `Resolve the paths in profile
	`,
	Run: runResolve,
}

func init() {
	RootCmd.AddCommand(resolveCmd)

	buildCmd.Flags().BoolVarP(&resolve_save, "save", "s", false, "If the Profile should be Updated with the new files")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runResolve(cmd *cobra.Command, args []string) {
	profilePath, err := filepath.Abs(profile)
	if err != nil {
		fmt.Printf("Cannot resolve absolute path of %s: %s", profile, err)
		os.Exit(1)
	}
	profile, err := image.LoadProfile(profilePath)
	if err != nil {
		fmt.Printf("%s is not a Profile: %s", profilePath, err)
		os.Exit(1)
	}
	config,err := image.LoadConfiguration(cfgFile)
	if err != nil {
		fmt.Printf("Cannot Load Configuration: %s", err)
		os.Exit(1)
	}
	files := config.GetFiles(profile.FileSets)
	for _, file := range files {
		fmt.Println(file)
	}
	if resolve_save {
		profile.Files = files
		err := profile.Save(filepath.Dir(profilePath))
		if err != nil {
			fmt.Printf("Cannot Save Profile: %s\n", err)
		}
	}
	if err != nil {
		fmt.Printf("Image Creation Failed: %s\n", err)
		os.Exit(1)
	}
}