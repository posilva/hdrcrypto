/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"hdrcrypto/pkg/hedera"
	"os"
	"path/filepath"
)

const (
	defaultConfigFilename = "hdrcrypto"
	defaultConfigPath     = "."
	defaultConfigType     = "yml"
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
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

type ConfigurationYaml struct {
	TokenId string
}

var hdrClient *hedera.HDRClient
var viperConfig *viper.Viper
var AppConfig ConfigurationYaml

func init() {
	viperConfig = viper.New()

	// set the configs defaults

	viperConfig.SetDefault("address", ":3000")

	// setup config file
	viperConfig.SetConfigName(defaultConfigFilename)

	viperConfig.SetConfigType(defaultConfigType)
	viperConfig.AddConfigPath(defaultConfigPath)

	// read configuration
	if err := viperConfig.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Errorf("Config file %s not found ", viperConfig.ConfigFileUsed())
		}
	}

	// setup environment variables configuration
	viperConfig.SetEnvPrefix(envPrefix)
	viperConfig.AutomaticEnv()

	err := viperConfig.Unmarshal(&AppConfig)
	if err != nil {
		panic(err)
	}

	// watching for config changes
	viperConfig.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	viperConfig.WatchConfig()

	setupClient()
}

func setupClient() {
	hdrClient = hedera.NewClientForTestNet()

	opId := viperConfig.GetString("operator_id")
	opKey := viperConfig.GetString("operator_key")

	err := hdrClient.Operator(opId, opKey)
	if err != nil {
		panic(err)
	}
}

func ConfigFileNamePath() string {
	return filepath.Join(defaultConfigPath, defaultConfigFilename) + "." + defaultConfigType
}
