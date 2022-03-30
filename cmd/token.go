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
	"github.com/spf13/cobra"
	"hdrcrypto/pkg/hedera"
)

// tokenCmd represents the token command
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Allow to interact with token service",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if viperConfig.Get("operator_id") == "" && viperConfig.Get("operator_key") == "" {
			return errors.New("Required operator account id and/or private key not set ")
		}
		setupClient()
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var hdrClient *hedera.HDRClient

func setupClient() {
	hdrClient = hedera.NewClientForTestNet()

	opId := viperConfig.GetString("operator_id")
	opKey := viperConfig.GetString("operator_key")

	err := hdrClient.Operator(opId, opKey)
	if err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.AddCommand(tokenCmd)
	tokenCmd.PersistentFlags().BoolP("force", "f", false, "Force to execute the action")
	viperConfig.BindPFlag("force", tokenCmd.PersistentFlags().Lookup("force"))

	tokenCmd.Flags().String("operator_id", "", "Operator account id")
	tokenCmd.Flags().String("operator_key", "", "Operator account private key")

	_ = viperConfig.BindPFlag("operator_id", tokenCmd.Flags().Lookup("operator_id"))
	_ = viperConfig.BindPFlag("operator_key", tokenCmd.Flags().Lookup("operator_key"))
	_ = viperConfig.BindEnv("operator_id")
	_ = viperConfig.BindEnv("operator_key")
}
