/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/base64"
	"fmt"

	"github.com/benduncan/poh-golang/pkg/wallet"
	"github.com/spf13/cobra"
)

// pubkeyCmd represents the pubkey command
var pubkeyCmd = &cobra.Command{
	Use:   "pubkey",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")

		mywallet, _ := wallet.Load(path)

		fmt.Printf("%s\n", base64.StdEncoding.EncodeToString(mywallet.PublicKey))

	},
}

func init() {
	keygenCmd.AddCommand(pubkeyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pubkeyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pubkeyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	pubkeyCmd.Flags().StringP("path", "p", ".wallet.json", "Wallet filename")

}
