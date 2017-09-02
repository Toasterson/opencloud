package cmd

import (
	"fmt"

	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/toasterson/opencloud/image"
)

var (
	build_type       string
	build_notresolve bool
	build_output     string
	build_profile    string
)

const build_default_baseDir string = "/tmp"

// buildCmd represents the create command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build an image",
	Long: `Build an image use -t to overwrite the type
	`,
	Run: runBuild,
}

func init() {
	RootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringVarP(&build_type, "type", "t", "", "use to overwrite the type of image to build")
	buildCmd.Flags().BoolVarP(&build_notresolve, "notresolve", "n", false, "Use files in profile and not resolve from Filesystem")
	buildCmd.Flags().StringVarP(&build_output, "output", "o", build_default_baseDir, "The Directory where a build will be saved to. defaults to /tmp/imgname")
	buildCmd.Flags().StringVarP(&build_profile, "profile", "p", "./profile.json", "The Profile file to use. Defaults to profile.json in pwd")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runBuild(cmd *cobra.Command, args []string) {
	profilePath, err := filepath.Abs(build_profile)
	if err != nil {
		fmt.Printf("Cannot resolve absolute path of %s: %s", build_profile, err)
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
	if !notresolve {
		profile.Files = config.GetFiles(profile.FileSets)
	}
	if build_type == "" {
		build_type = string(profile.Type)
	}
	switch build_type {
	case image.TypeChroot:
		err = image.BuildChroot(profile, build_output)
	case image.TypeACI:
		err = image.BuildACI(profile, build_output)
	case image.TypeTar:
		err = image.BuildTar(profile, build_output)
	case image.TypeUfs:
		err = image.BuildUFS(profile, build_output)
	case image.TypeZfs:
		err = image.BuildZFS(profile, build_output)
	default:
		fmt.Printf("Format %s not known", build_type)
		os.Exit(1)
	}
	if err != nil {
		fmt.Printf("Image Creation Failed: %s", err)
		os.Exit(1)
	}
}