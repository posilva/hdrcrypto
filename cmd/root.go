/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

const (
	defaultConfigFilename = ".hdrcrypto"
	envPrefix             = "HDRCRYPTO"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hdrcrypto",
	Short: "Command line tool to demonstrate the Hedera Token operation",
	Long:  `Command line tool to demonstrate the Hedera Token operation`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var v *viper.Viper

func init() {
	v = viper.New()

	// set the configs defaults
	v.SetDefault("address", ":3000")

	// setup config file
	v.SetConfigName(defaultConfigFilename)
	v.SetConfigType("yaml")
	v.AddConfigPath("$HOME/.hdrcrypto")
	v.AddConfigPath(".")

	// read configuration
	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Errorf("Config file %s not found ", v.ConfigFileUsed())
		}
	}

	// setup environment variables configuration
	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()

	// watching for config changes
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	v.WatchConfig()
}
