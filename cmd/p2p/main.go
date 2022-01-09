package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/benduncan/poh-golang/pkg/p2pnet"
)

func main() {

	//poh := poh_hash.New("")
	p2p := p2pnet.New()
	//p2p.POH = &poh

	var server = flag.Int("s", 0, "Launch server")

	flag.Parse()

	if *server > 0 {

		fmt.Println("Launching server")

		go func() {
			p2p.BroadcastListen(p2p.MsgHandler)
		}()

		go func() {
			//poh.GeneratePOH(100_000_000_000)
		}()

		for {

		}

	} else {

		fmt.Println("Launching client")

		for i := 0; i < 1_000; i++ {
			p2p.BroadcastSend(fmt.Sprintf("%d", i))
			time.Sleep(1 * time.Millisecond)
		}

	}

	/*
	 */

}
