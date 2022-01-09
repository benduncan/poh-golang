/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/benduncan/poh-golang/pkg/p2pnet"
	"github.com/benduncan/poh-golang/pkg/poh_hash"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

const walletLocation = "path"
const flagIP = "ip"
const flagPort = "port"

type POH struct {
	record poh_hash.POH
}

type SyncState struct {
	Entry []poh_hash.POH_Entry
	Len   int
}

var poh POH

func (this *POH) syncstate(c *gin.Context) {

	//this.Mu.RLock()

	fmt.Println("Length =>", len(this.record.POH[0].Entry))
	fmt.Println(this.record.POH[0].Entry)

	var records []poh_hash.POH_Entry
	var len = len(this.record.POH[0].Entry) - 1

	start := 1
	if len > 10 {
		start = len - 10
	}

	records = append(records, this.record.POH[0].Entry[0])

	for a := start; a < len; a++ {
		records = append(records, this.record.POH[0].Entry[a])
	}
	c.JSON(200, SyncState{Entry: records, Len: len})

	//poh.Mu.RUnlock()

}

func (this *POH) syncdatastate(c *gin.Context) {

	// TODO: Find more efficient way to handle returning entries, gin-tonic limitation it seems for c.JSON
	c.Data(200, "application/json; charset=utf-8", []byte(fmt.Sprintf("{\"PublicKey\": \"%s\", \"Data\": [", base64.StdEncoding.EncodeToString(this.record.Wallet.PublicKey))))

	//c.Data(200, "application/json; charset=utf-8", []byte("["))

	// Print the first
	this.record.Mu.RLock()
	c.JSON(200, this.record.POH[0].Entry[0])

	len := len(this.record.POH[0].Entry)

	for a := 1; a < len; a++ {

		if this.record.POH[0].Entry[a].Data != nil {
			c.Data(200, "application/json; charset=utf-8", []byte(","))

			c.JSON(200, this.record.POH[0].Entry[a])

		}

	}

	this.record.Mu.RUnlock()

	c.Data(200, "application/json; charset=utf-8", []byte("]}"))

}

func (this *POH) verify(c *gin.Context) {

	host, _ := c.GetQuery("host")

	resp, err := http.Get(fmt.Sprintf("http://%s:8080/syncdata", host))

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	var valid struct {
		PublicKey string
		Data      []poh_hash.POH_Entry
	}

	validator := valid

	json.NewDecoder(resp.Body).Decode(&validator)

	confirmation := poh_hash.POH{}
	pubkey, _ := base64.StdEncoding.DecodeString(validator.PublicKey)

	confirmation.Wallet.PublicKey = pubkey

	confirmation.POH = append(confirmation.POH, poh_hash.POH_Epoch{Epoch: 1})
	confirmation.POH[0].Entry = validator.Data

	err = confirmation.VerifyPOH(8)

	if err == nil {
		c.JSON(200, gin.H{"Status": "OK"})
	} else {
		c.JSON(500, gin.H{"Status": "fail"})
		fmt.Println("VerifyPOH error =>", err)
	}

}

func (this *POH) pushstate(c *gin.Context) {

	//this.Mu.RLock()

	data, _ := c.GetQuery("data")
	sender, _ := c.GetQuery("sender")

	queuedata := poh_hash.QueueData{Data: []byte(data), Sender: []byte(sender)}
	this.record.QueueSync.State = append(this.record.QueueSync.State, queuedata)

	c.JSON(200, queuedata)

	//c.JSON(200, SyncState{Entry: records, Len: len})

	//poh.Mu.RUnlock()

}

func (this *POH) state(c *gin.Context) {

	c.JSON(200, this.record.QueueSync)

}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ip, _ := cmd.Flags().GetString(flagIP)
		port, _ := cmd.Flags().GetUint64(flagPort)
		walletPath, _ := cmd.Flags().GetString(walletLocation)

		// Check wallet path
		if _, err := os.Stat(walletPath); err != nil {
			log.Fatalf("Wallet %s could not be opened (%s)", walletPath, err)
		}

		fmt.Println(ip, port)

		// Creates a gin router with default middleware:
		// logger and recovery (crash-free) middleware
		router := gin.Default()

		poh.record = poh_hash.New(walletPath)

		fmt.Println(poh)

		p2p := p2pnet.New()
		p2p.POH = &poh.record

		// Launch the UDP packet receiver
		go func() {
			p2p.BroadcastListen(p2p.MsgHandler)
		}()

		// Launch the PoH go routine
		go func() {
			poh.record.GeneratePOH(100_000_000_000)
		}()

		fmt.Println(poh)

		router.GET("/sync", poh.syncstate)

		router.GET("/syncdata", poh.syncdatastate)

		router.GET("/verify", poh.verify)

		router.GET("/push", poh.pushstate)

		router.GET("/state", poh.state)

		// By default it serves on :8080 unless a
		// PORT environment variable was defined.
		router.Run()

	},
}

func init() {

	rootCmd.PersistentFlags().String(walletLocation, ".wallet.json", "Filename for wallet location to launch peer")
	rootCmd.PersistentFlags().String(flagIP, "127.0.0.1", "exposed IP for communication with peers")
	rootCmd.PersistentFlags().Uint64(flagPort, 24816, "exposed HTTP port for communication with peers")
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
