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
	"hdrcrypto/pkg/hedera"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// tokenCmd represents the token command
var accountCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Allow to create an account",
	Long:  `Allow to create an account`,
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("name", cmd.Flags().Lookup("name"))
		viper.BindPFlag("balance", cmd.Flags().Lookup("balance"))

	},
	Run: func(cmd *cobra.Command, args []string) {
		account, err := hedera.CreateAccountEntity(hdrClient, viper.GetString("name"), viper.GetFloat64("balance"))
		if err != nil {
			panic(err)
		}
		log.Info().Msgf("created account: %v", account)
	},
}

func init() {
	accountCmd.AddCommand(accountCreateCmd)
	accountCreateCmd.Flags().StringP("name", "n", "", "Name of the account")
	accountCreateCmd.MarkFlagRequired("name")
	accountCreateCmd.Flags().Float64P("balance", "b", 1, "Initial balance of hbar")
	accountCreateCmd.MarkFlagRequired("balance")
}
