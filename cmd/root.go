// Package cmd
/*
Copyright © 2022 Pedro Marques da Silva <posilva@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"hdrcrypto/pkg/hederalib"
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

var hdrClient *hederalib.HDRClient
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
	hdrClient = hederalib.NewClientForTestNet()

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
