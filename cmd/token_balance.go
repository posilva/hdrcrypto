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
	"github.com/hashgraph/hedera-sdk-go/v2"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tokenBalanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Command to get the Fungible Token balance to a given token.",
	Long:  `Command to get the Fungible Token balance to a given token.`,
	Run: func(cmd *cobra.Command, args []string) {
		account := hdrClient.OperatorAccount()

		if viper.IsSet("account") {
			a, err := hedera.AccountIDFromString(viper.GetString("account"))
			if err != nil {
				panic(err)
			}
			account = a
		}
		query := hedera.NewAccountBalanceQuery().
			SetAccountID(account)

		accountBalance, err := query.Execute(hdrClient.Get())
		if err != nil {
			panic(err)
		}

		log.Info().Msgf("The HDR token account balance for account '%v' is %v", account, accountBalance.Tokens.Get(activeToken.Id))
		log.Info().Msgf("The HBAR account balance for this account '%v' is %v", account, accountBalance.Hbars.String())
	},
}

func init() {
	tokenCmd.AddCommand(tokenBalanceCmd)
	tokenBalanceCmd.Flags().StringP("account", "a", "", "Target Account Id")
	viper.BindPFlag("account", tokenBalanceCmd.Flags().Lookup("account"))
}
