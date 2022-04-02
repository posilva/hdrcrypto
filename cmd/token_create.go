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
)

var tokenCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Command to create a Fungible Token.",
	Long:  `Command to create a Fungible Token.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := hederalib.TokenConfig{
			Name:                  "HDR Crypto",
			Symbol:                "HDR",
			MaxTransactionFeeHbar: 30,
			InitialSupply:         5000,
			Decimals:              2,
		}
		token, err := hederalib.NewToken(hdrClient, config)
		if err != nil {
			panic(err)
		}

		log.Info().Msgf("Created new token: %s", token)
	},
}

func init() {
	tokenCmd.AddCommand(tokenCreateCmd)
}
