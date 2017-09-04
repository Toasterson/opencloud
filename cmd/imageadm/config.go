package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/toasterson/opencloud/image"
)

// createCmd represents the create command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Shows the Configuration",
	Long: `Shows the Current configuration
	`,
	Run: configCmdrun,
}

func init() {
	RootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func configCmdrun(cmd *cobra.Command, args []string) {
	conf, err := image.LoadConfiguration(cfgFile)
	if err != nil {
		panic(err)
	}
	fmt.Println("Configuration:")
	for _, section := range conf.Sections{
		fmt.Printf("[%s]\n", section.Name)
		if section.Comment != ""{
			fmt.Println(section.Comment)
		}
		if len(section.Paths) > 0 {
			fmt.Printf("Paths: %v\n", section.Paths)
		}
		if len(section.Dependencies) > 0 {
			fmt.Printf("Dependencies: %v\n", section.Dependencies)
		}
		if len(section.Users) > 0 {
			fmt.Printf("Users: %v\n", section.Users)
		}
		if len(section.Groups) > 0 {
			fmt.Printf("Groups: %v\n", section.Groups)
		}
		if len(section.Devices) > 0 {
			fmt.Printf("Devices: %v\n", section.Devices)
		}
		fmt.Println()
	}
}
