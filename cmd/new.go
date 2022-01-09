/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/benduncan/poh-golang/pkg/wallet"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		force, _ := cmd.Flags().GetBool("force")
		path, _ := cmd.Flags().GetString("path")
		name, _ := cmd.Flags().GetString("name")

		if path == "" {
			path = ".wallet.json"
		}

		mywallet := wallet.New()
		mywallet.Name = name

		err := mywallet.GenerateWallet()

		if err != nil {
			log.Fatal(err)
		}

		err = mywallet.Save(path, force)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("New wallet saved in %s\n", path)
		fmt.Println("Your wallet public key is:")
		fmt.Printf("%s\n", base64.StdEncoding.Strict().EncodeToString(mywallet.PublicKey))

	},
}

func init() {
	keygenCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	newCmd.Flags().BoolP("force", "f", false, "Force overwrite")
	newCmd.Flags().StringP("path", "p", ".wallet.json", "Wallet filename")
	newCmd.Flags().StringP("name", "n", "My Wallet", "Wallet identifier name (My wallet)")

}
