package main

import (
	"fmt"
	"runtime"

	"github.com/benduncan/poh-golang/pkg/poh_hash"
)

func main() {

	fmt.Println("Proof of History example")
	cpu_cores := runtime.NumCPU()
	fmt.Printf("CPU Cores: %d\n", cpu_cores)

	poh := poh_hash.New()
	poh.GeneratePOH(10_000_000)

	fmt.Printf("Generate Hashrate %d p/sec (1-core)\n", poh.HashRate)

	poh.VerifyPOH(cpu_cores)

	fmt.Printf("Verify Hashrate %d p/sec (%d-cores)\n", poh.VerifyHashRate, cpu_cores)
	fmt.Printf("Verify Hashrate %d p/core\n", poh.VerifyHashRatePerCore)

}
