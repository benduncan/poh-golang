package main

import (
	"fmt"
	"runtime"

	"github.com/benduncan/poh-golang/pkg/poh_hash"
)

func main() {

	var count = 20_000_000
	fmt.Println("Proof of History example hash-generation and verification.")

	cpu_cores := runtime.NumCPU()
	fmt.Printf("CPU Cores: %d\n", cpu_cores)

	fmt.Printf("Generating %d hashes\n", count)
	poh := poh_hash.New("")
	poh.GeneratePOH(uint64(count))

	fmt.Printf("Generate Hashrate %d p/sec (1-core)\n", poh.HashRate)

	poh.VerifyPOH(cpu_cores)

	fmt.Printf("Verify Hashrate %d p/sec (%d-cores)\n", poh.VerifyHashRate, cpu_cores)
	fmt.Printf("Verify Hashrate %d p/core\n", poh.VerifyHashRatePerCore)
	fmt.Printf("Verification step %.2fx vs single core hash generation", float64(poh.VerifyHashRate)/float64(poh.HashRate))

}
