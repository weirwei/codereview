package cmd

import (
	"fmt"
	"os/user"

	"github.com/spf13/cobra"
	"github.com/weirwei/codereview/log"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Config CodeReview",
	Long:  `Config CodeReview`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var setCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]
		userViper.Set(key, value)
		if err := userViper.WriteConfig(); err != nil {
			fmt.Printf("Error writing config: %v\n", err)
			return
		}
		fmt.Printf("Set config '%s' to '%s'\n", key, value)
	},
}

var getCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := userViper.Get(key)
		switch value := value.(type) {
		case string:
			fmt.Printf("%s=%v\n", key, value)
		case map[string]interface{}:
			for k, v := range value {
				fmt.Printf("%s.%s=%v\n", key, k, v)
			}
		default:
			fmt.Printf("%s=%v\n", key, value)
		}
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration values",
	Run: func(cmd *cobra.Command, args []string) {
		allSettings := userViper.AllSettings()
		for key, value := range allSettings {
			switch v := value.(type) {
			case string:
				fmt.Printf("%s=%v\n", key, v)
			case map[string]interface{}:
				for subKey, subValue := range v {
					fmt.Printf("%s.%s=%v\n", key, subKey, subValue)
				}
			default:
				fmt.Printf("%s=%v\n", key, v)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	// viper.SetConfigType("ini")
	// home := GetUserHomeDir()
	// viper.SetConfigFile(path.Join(home, ".codereview"))
	// err := viper.ReadInConfig()
	// if _, ok := err.(viper.ConfigFileNotFoundError); ok {
	// 	err := viper.WriteConfigAs(viper.ConfigFileUsed())
	// 	if err != nil {
	// 		log.Fatalf("创建配置文件失败: %s", err)
	// 		return
	// 	}
	// }
	configCmd.AddCommand(setCmd, getCmd, listCmd)
}

func GetUserHomeDir() string {
	// Unix系统
	u, err := user.Current()
	if err != nil {
		log.Error(err)
		return ""
	}

	return u.HomeDir
}
