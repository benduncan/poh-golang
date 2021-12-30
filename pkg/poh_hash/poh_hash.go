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

type POH_Entry struct {
	Hash []byte
	Data []byte
	Seq  uint64
}

type POH_Epoch struct {
	Entry []POH_Entry
	Epoch uint32
}

type POH struct {
	POH                   map[uint32]POH_Epoch
	HashRate              uint32
	VerifyHashRate        uint32
	VerifyHashRatePerCore uint32
	TickRate              uint32
	mu                    sync.RWMutex
}

func New() (this POH) {

	// todo confirm make len vs append
	this.POH = make(map[uint32]POH_Epoch)

	this.TickRate = 10

	return

}

func (this *POH) GeneratePOH(count int) {

	start := time.Now()

	tickcount := count / int(this.TickRate)
	this.POH[0] = POH_Epoch{Epoch: 1, Entry: make([]POH_Entry, tickcount)}

	h := sha256.New()
	h.Write([]byte(genisis_hash))

	this.mu.Lock()
	this.POH[0].Entry[0] = POH_Entry{Hash: h.Sum(nil), Seq: 0}
	this.mu.Unlock()

	var prevhash []byte

	for i := uint64(1); i < uint64(count); i++ {
		h := sha256.New()

		t := i % uint64(this.TickRate)

		// Only save a state every X events (based on TickSize) to reduce memory allocation
		if t == 0 {
			a := i / uint64(this.TickRate)
			h.Write([]byte(hex.EncodeToString(this.POH[0].Entry[a-1].Hash)))
			hash := h.Sum(nil)
			this.mu.Lock()
			this.POH[0].Entry[a] = POH_Entry{Hash: hash, Seq: i}
			this.mu.Unlock()
		} else {
			h.Write([]byte(hex.EncodeToString(prevhash)))
			prevhash = h.Sum(nil)

		}

	}

	timer := time.Now()
	elapsed := timer.Sub(start)

	//fmt.Printf("Poh (%d) loops processed in (%s)\n", count, elapsed)
	//fmt.Printf("Loops processed in secs (%f) (%f)\n", elapsed.Seconds(), (1 / elapsed.Seconds()))

	this.HashRate = uint32(float64(count) * (1 / elapsed.Seconds()))

	return
}

func (this *POH) VerifyPOH(cpu_cores int) {

	start := time.Now()

	pool := pond.New(cpu_cores, 0, pond.MinWorkers(cpu_cores))

	// Distribute jobs on each core for the specified task-size
	tasksize := uint64(100_000)
	tasks := uint64(len(this.POH[0].Entry)) / tasksize

	for i := uint64(0); i < tasks; i++ {
		n := i
		pool.Submit(func() {

			seqstart := n * tasksize
			seqend := seqstart + tasksize

			var prevhash []byte

			for i := seqstart; i < seqend; i++ {

				h := sha256.New()

				t := i % uint64(this.TickRate)

				if t == 0 {

					var a uint64

					// Confirm if the genisis block
					if i == 0 {
						h.Write([]byte(genisis_hash))
					} else {
						// Find the index of the last previous state 'save' from the TickRate
						a = i / uint64(this.TickRate)
						h.Write([]byte(hex.EncodeToString(this.POH[0].Entry[a-1].Hash)))
					}

					// Read lock
					this.mu.RLock()
					orig := this.POH[0].Entry[a].Hash
					this.mu.RUnlock()

					// Compare the proof to the original
					compare := bytes.Compare(h.Sum(nil), orig)

					if compare == 1 {
						fmt.Printf("Seq Hash %d => %s\n", i, hex.EncodeToString(h.Sum(nil)))
						fmt.Printf("Orig Hash %d => %s\n", i, hex.EncodeToString(orig))
						log.Fatalf("Error! Proof does not match")
					}

				} else {
					h.Write([]byte(hex.EncodeToString(prevhash)))
					prevhash = h.Sum(nil)
				}

			}

		})
	}

	// Stop the pool and wait for all submitted tasks to complete
	pool.StopAndWait()

	timer := time.Now()
	elapsed := timer.Sub(start)

	// Calculate the verification hashrate
	this.VerifyHashRate = uint32(float64(len(this.POH[0].Entry)) * (1 / elapsed.Seconds()))
	this.VerifyHashRatePerCore = this.VerifyHashRate / uint32(cpu_cores)

}
