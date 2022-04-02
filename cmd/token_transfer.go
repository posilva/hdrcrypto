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
FITNESS FOR A PARTICULAR PURPOSE AND NON INFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"hdrcrypto/pkg/hederalib"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tokenTransferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Command to trasfer a Fungible Token.",
	Long:  `Command to transfer a Fungible Token.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("amount-value", cmd.Flags().Lookup("amount-value"))
		viper.BindPFlag("account", cmd.Flags().Lookup("account"))
	},
	Run: func(cmd *cobra.Command, args []string) {

		err := hederalib.TokenTransfer(hdrClient,
			activeToken.Id.String(),
			viper.GetString("account"),
			viper.GetFloat64("amount-value"))

		if err != nil {
			log.Error().Msgf("%v", err)
		}
		log.Logger.Info().Msgf("processing transfer: %v", args)
	},
}

func init() {
	tokenCmd.AddCommand(tokenTransferCmd)
	tokenTransferCmd.Flags().StringP("account", "a", "", "Target account to transfer token ammount to")
	tokenTransferCmd.MarkFlagRequired("account")

	tokenTransferCmd.Flags().StringP("amount-value", "v", "", "Amount to be transfered to target account")
	tokenTransferCmd.MarkFlagRequired("amount-value")

}
