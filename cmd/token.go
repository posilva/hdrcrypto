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
	"errors"
	"hdrcrypto/pkg/hederalib"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var activeToken *hederalib.Token

// tokenCmd represents the token command
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Allow to interact with token service",
	Long:  `Allow to interact with token service`,
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("operator_id", cmd.Flags().Lookup("operator_id"))
		viper.BindPFlag("operator_key", cmd.Flags().Lookup("operator_key"))
		viper.BindEnv("operator_id")
		viper.BindEnv("operator_key")
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if viper.Get("operator_id") == "" && viper.Get("operator_key") == "" {
			return errors.New("required operator account id and/or private key not set ")
		}

		loadTokenFromConfig()
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if activeToken == nil {
			log.Error().Msgf("missing token id in configuration file")
		}
	},
}

func loadTokenFromConfig() {
	if AppConfig.TokenId != "" {
		token, err := hederalib.NewTokenFromInfo(hdrClient, AppConfig.TokenId)
		if err != nil {
			panic(err)
		}
		activeToken = token
		return
	}
}

func init() {
	rootCmd.AddCommand(tokenCmd)
	tokenCmd.Flags().String("operator_id", "", "Operator account id")
	tokenCmd.Flags().String("operator_key", "", "Operator account private key")
}
