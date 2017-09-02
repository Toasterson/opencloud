package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy an Image to the following Opencloud Realm",
	Long: `Deploys an image to an opencloud Realm.

	use -h to specify another Opencloud Installation than the public one.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)

	createCmd.Flags().StringP("host", "H", "http://opencloud.wegmueller.it", "The Url to publish to.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
