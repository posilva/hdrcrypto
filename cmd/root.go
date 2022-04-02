/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"hdrcrypto/pkg/hedera"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ConfigurationYaml struct {
	TokenId string
}

const (
	defaultConfigFilename = "hdrcrypto"
	defaultConfigPath     = "."
	defaultConfigType     = "yml"
	envPrefix             = "HDRCRYPTO"
)

var hdrClient *hedera.HDRClient
var AppConfig ConfigurationYaml

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

func init() {

	viper.SetDefault("address", ":3000")

	// setup config file
	viper.SetConfigName(defaultConfigFilename)

	viper.SetConfigType(defaultConfigType)
	viper.AddConfigPath(defaultConfigPath)

	// read configuration
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Error().Msgf("Config file %v not found ", viper.ConfigFileUsed())
		}
	}

	// setup environment variables configuration
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()

	err := viper.Unmarshal(&AppConfig)
	if err != nil {
		panic(err)
	}

	// watching for config changes
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info().Msgf("Config file changed:", e.Name)
	})
	viper.WatchConfig()
	setupClient()
}

func setupClient() {
	hdrClient = hedera.NewClientForTestNet()

	opId := viper.GetString("operator_id")
	opKey := viper.GetString("operator_key")

	if opId == "" || opKey == "" {
		log.Fatal().Msgf("required parameters operator_id or operator_key missing")
		return

	}
	err := hdrClient.Operator(opId, opKey)
	if err != nil {
		panic(err)
	}
}

func ConfigFileNamePath() string {
	return filepath.Join(defaultConfigPath, defaultConfigFilename) + "." + defaultConfigType
}
