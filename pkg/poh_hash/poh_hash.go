package poh_hash

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/alitto/pond"
)

const genisis_hash = "GENISIS_STRING"

type POH struct {
	Hash                  [][]byte
	mu                    sync.RWMutex
	HashRate              int64
	VerifyHashRate        int64
	VerifyHashRatePerCore int64
}

func New() (poh POH) {
	return

}

func (poh *POH) GeneratePOH(count int) {

	start := time.Now()
	poh.Hash = make([][]byte, count)

	h := sha256.New()
	h.Write([]byte(genisis_hash))

	poh.mu.Lock()
	poh.Hash[0] = h.Sum(nil)
	poh.mu.Unlock()

	//fmt.Println("Genesis Hash ", hex.EncodeToString(poh.Hash[0]))

	for i := 1; i < count; i++ {
		h := sha256.New()
		h.Write([]byte(hex.EncodeToString(poh.Hash[i-1])))
		poh.mu.Lock()
		poh.Hash[i] = h.Sum(nil)
		poh.mu.Unlock()
	}

	timer := time.Now()
	elapsed := timer.Sub(start)

	//fmt.Printf("Poh (%d) loops processed in (%s)\n", count, elapsed)
	//fmt.Printf("Loops processed in secs (%f) (%f)\n", elapsed.Seconds(), (1 / elapsed.Seconds()))

	poh.HashRate = int64(float64(count) * (1 / elapsed.Seconds()))

	return

}

func (poh *POH) VerifyPOH(cpu_cores int) {

	start := time.Now()

	pool := pond.New(cpu_cores, 0, pond.MinWorkers(cpu_cores))

	tasksize := 1_000_000
	tasks := len(poh.Hash) / tasksize

	// Submit 1000 tasks
	for i := 0; i < tasks; i++ {
		n := i
		pool.Submit(func() {

			//block := count / 20
			seqstart := n * tasksize
			seqend := seqstart + tasksize

			//fmt.Println(n, tasksize, seqstart, seqend)

			for i := seqstart; i < seqend; i++ {

				var validate []byte
				h := sha256.New()

				if i == 0 {
					validate = []byte(genisis_hash)
					h.Write(validate)
				} else {
					validate = poh.Hash[i-1]
					h.Write([]byte(hex.EncodeToString(validate)))
				}

				proof := h.Sum(nil)

				poh.mu.RLock()
				orig := &poh.Hash[i]
				poh.mu.RUnlock()

				compare := bytes.Compare(proof, *orig)

				if compare == 0 {
					//fmt.Println("Proof match")
				} else {
					log.Fatalf("Error! Proof does not match")
					fmt.Printf("Seq Hash %d => %s\n", i, hex.EncodeToString(proof))
					fmt.Printf("Orig Hash %d => %s\n", i, hex.EncodeToString(poh.Hash[i]))

				}

			}

			//fmt.Printf("Running task #%d\n", n)
		})
	}

	// Stop the pool and wait for all submitted tasks to complete
	pool.StopAndWait()

	timer := time.Now()
	elapsed := timer.Sub(start)

	poh.VerifyHashRate = int64(float64(len(poh.Hash)) * (1 / elapsed.Seconds()))
	poh.VerifyHashRatePerCore = poh.VerifyHashRate / int64(cpu_cores)

}
