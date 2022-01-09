/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"runtime"

	"github.com/benduncan/poh-golang/pkg/poh_hash"
	"github.com/spf13/cobra"
)

// localCmd represents the local command
var localCmd = &cobra.Command{
	Use:   "local",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		walletPath, _ := cmd.Flags().GetString(walletLocation)

		fmt.Println("Proof of History example")

		cpu_cores := runtime.NumCPU()
		fmt.Printf("CPU Cores: %d\n", cpu_cores)

		poh := poh_hash.New(walletPath)
		poh.GeneratePOH(100_000_000)

		fmt.Printf("Generate Hashrate %d p/sec (1-core)\n", poh.HashRate)

		poh.VerifyPOH(cpu_cores)

		fmt.Printf("Verify Hashrate %d p/sec (%d-cores)\n", poh.VerifyHashRate, cpu_cores)
		fmt.Printf("Verify Hashrate %d p/core\n", poh.VerifyHashRatePerCore)

	},
}

func init() {
	rootCmd.AddCommand(localCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// localCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// localCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
