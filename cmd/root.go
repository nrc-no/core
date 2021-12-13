package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFiles []string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "core",
	Short: "Data collection and case management for humanitarian organizations",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringSliceVar(&cfgFiles, "config", []string{}, "config files")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	viper.AutomaticEnv() // read in environment variables that match

	for i, file := range cfgFiles {
		viper.SetConfigFile(file)
		var err error
		if i == 0 {
			err = viper.ReadInConfig()
		} else {
			err = viper.MergeInConfig()
		}
		if err == nil {
			fmt.Println("Using config file: " + viper.ConfigFileUsed())
		} else {
			fmt.Println(err.Error())
			panic(err)
		}
	}

	if len(cfgFiles) > 0 {
		viper.WatchConfig()
	}

}
