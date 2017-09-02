package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/toasterson/opencloud/image"
)

var (
	add bool
	remove bool
)

// createCmd represents the create command
var setConfigCmd = &cobra.Command{
	Use:   "set-config",
	Short: "Changes the Config",
	Long: `Change the configuration from cmdline
	use -s to specify which section to edit
	use -a to add to existing and -r to remove one entry only
	use name=value as positional arguments to specify what to change. multiple positional arguments are accepted
	`,
	Run: setConfigCmdrun,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("At least one argument required")
		}
		for _, arg := range args{
			if !strings.Contains(arg, "="){
				return fmt.Errorf("Arguments must have name=value syntax. %s does not", arg)
			}
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(setConfigCmd)

	setConfigCmd.Flags().StringP("section", "s", "", "The Section to edit")
	setConfigCmd.Flags().BoolVarP(&add, "add", "a", false, "Add the value to the section instead of overwriting it")
	setConfigCmd.Flags().BoolVarP(&remove,"remove", "r", false, "Remove the value from the section")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func setConfigCmdrun(cmd *cobra.Command, args []string) {
	section := cmd.Flags().Lookup("section").Value.String()
	conf, err := image.LoadConfiguration(cfgFile)
	if err != nil {
		conf = image.Config{}
	}
	sectionObj, ok := conf.Sections[section]
	if !ok {
		sectionObj = image.ConfigSection{Name: section}
	}
	if remove{
		arg := args[0]
		if strings.Contains(arg, "="){
			tmp := strings.SplitN(arg, "=", -1)
			name := tmp[0]
			val := tmp[1]
			removeEntry(&sectionObj, name, val)
		} else {
			switch arg {
			case "devices", "dev":
				sectionObj.Devices = []string{}
			case "users", "user":
				sectionObj.Users = []string{}
			case "groups", "group":
				sectionObj.Groups = []string{}
			case "comment":
				sectionObj.Comment = ""
			case "paths":
				sectionObj.Paths = []string{}
			case "deps", "dependencies":
				sectionObj.Dependencies = []string{}
			default:
				panic(fmt.Errorf("%s does not exist", arg))
			}
		}
	} else {
		argMap := argsToMap(args)
		for name, values := range argMap {
			apply(&sectionObj, name, values)
		}
	}
	err = image.SaveConfigSection(sectionObj)
	if err != nil {
		panic(err)
	}
}

func argsToMap(args []string) map[string][]string {
	argMap := make(map[string][]string)
	for _, arg := range args {
		tmp := strings.Split(arg, "=")
		name := tmp[0]
		val := tmp[1]
		_, ok := argMap[name]
		if ok {
			argMap[name] = append(argMap[name], val)
		} else {
			argMap[name] = []string{val}
		}
	}
	return argMap
}

func apply(section *image.ConfigSection, name string, values []string) {
	switch name {
	case "devices", "dev":
		if add {
			section.Devices = append(section.Devices, values...)
		} else {
			section.Devices = values
		}
	case "users", "user":
		if add {
			section.Users = append(section.Users, values...)
		} else {
			section.Users = values
		}
	case "groups", "group":
		if add {
			section.Groups = append(section.Groups, values...)
		} else {
			section.Groups = values
		}
	case "comment":
		section.Comment = values[0]
	case "paths":
		if add {
			section.Paths = append(section.Paths, values...)
		} else {
			section.Paths = values
		}
	case "deps", "dependencies":
		if add {
			section.Dependencies = append(section.Dependencies, values...)
		} else {
			section.Dependencies = values
		}
	default:
		panic(fmt.Errorf("%s does not exist", name))
	}
}

func removeEntry(section *image.ConfigSection, name string, value string) {
	switch name {
	case "devices", "dev":
		removeOneEntry(&section.Devices, value)
	case "users", "user":
		removeOneEntry(&section.Users, value)
	case "groups", "group":
		removeOneEntry(&section.Groups, value)
	case "comment":
		section.Comment = ""
	case "paths":
		removeOneEntry(&section.Paths, value)
	case "deps", "dependencies":
		removeOneEntry(&section.Dependencies, value)
	default:
		panic(fmt.Errorf("%s does not exist", name))
	}
}

func removeOneEntry(xs *[]string, entry string){
	j := 0
	for i, x := range *xs {
		if x != entry {
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]
}
