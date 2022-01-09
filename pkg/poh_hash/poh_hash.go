package poh_hash

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/alitto/pond"
	"github.com/benduncan/poh-golang/pkg/wallet"
)

var genisis_hash string

type POH_Entry struct {
	Hash      []byte
	Data      []byte
	Signature []byte
	Seq       uint64
}

type POH_Epoch struct {
	Entry []POH_Entry
	Epoch uint32
}

type POH struct {
	POH                   []POH_Epoch
	QueueSync             QueueSync
	HashRate              uint32
	VerifyHashRate        uint32
	VerifyHashRatePerCore uint32
	TickRate              uint64
	Mu                    sync.RWMutex
	Wallet                wallet.Wallet
}

type QueueSync struct {
	State []QueueData
}

type QueueData struct {
	Data   []byte
	Sender []byte
	Block  uint64
}

func New(walletpath string) (this POH) {

	this.TickRate = 1_000_000

	if walletpath == "" {
		walletpath = ".wallet.json"
	}

	var err error
	this.Wallet, err = wallet.Load(walletpath)

	if err != nil {
		err = this.Wallet.GenerateWallet()

		if err != nil {
			log.Fatalf("Could not create wallet %s", err)
		}

		err = this.Wallet.Save(walletpath, false)

		if err != nil {
			log.Fatalf("Could not save wallet file %s (%s)", walletpath, err)
		}

	}
	return

}

func (this *POH) FetchDataState(block uint64) (data []byte, chk bool) {

	for i := 0; i < len(this.QueueSync.State); i++ {

		if this.QueueSync.State[i].Block > 0 {
			continue
		} else {
			this.Mu.Lock()
			this.QueueSync.State[i].Block = block
			this.Mu.Unlock()
			return this.QueueSync.State[i].Data, true
		}

	}

	return data, false

}

func (this *POH) GeneratePOH(count uint64) {

	start := time.Now()

	//tickcount := count / this.TickRate
	this.POH = append(this.POH, POH_Epoch{Epoch: 1})

	h := sha256.New()
	var prevhash []byte

	// TODO: Geneisis + wallet public-key + timestamp + rand number
	genisis_hash := fmt.Sprintf("GENESIS_HASH-%d-%s", time.Now().UnixNano(), base64.StdEncoding.EncodeToString(this.Wallet.PrivateKey))

	/*
		d, _ := this.Wallet.Sign([]byte(genisis_hash))

		fmt.Println("Signed => ", string(d))

		fmt.Println(genisis_hash)

		verify, _ := this.Wallet.Verify([]byte("hello world"), d)

		fmt.Println("Verify => ", verify)
	*/

	h.Write([]byte(genisis_hash))
	prevhash = h.Sum(nil)

	this.Mu.Lock()
	this.POH[0].Entry = append(this.POH[0].Entry, POH_Entry{Hash: prevhash, Seq: 0})
	this.Mu.Unlock()

	for i := uint64(1); i < uint64(count); i++ {
		h := sha256.New()

		// TODO: Optimise, periodically push events published off the stack
		data, chk := this.FetchDataState(i)
		t := i % uint64(this.TickRate)

		// Create a new hash from a data request
		if chk == true {
			h.Write(append(prevhash, data...))
			prevhash = h.Sum(nil)

			signature, _ := this.Wallet.Sign(data)

			this.Mu.Lock()
			this.POH[0].Entry = append(this.POH[0].Entry, POH_Entry{Hash: prevhash, Data: data, Seq: i, Signature: signature})
			this.Mu.Unlock()

		} else {
			// Hash the latest output, hash of a hash for POH
			h.Write(prevhash)
			prevhash = h.Sum(nil)

		}

		if t == 0 {
			// Only save a state every X events (based on TickSize) to reduce memory allocation
			this.Mu.Lock()
			this.POH[0].Entry = append(this.POH[0].Entry, POH_Entry{Hash: prevhash, Seq: i})
			this.Mu.Unlock()

		}

	}

	timer := time.Now()
	elapsed := timer.Sub(start)

	//fmt.Printf("Poh (%d) loops processed in (%s)\n", count, elapsed)
	//fmt.Printf("Loops processed in secs (%f) (%f)\n", elapsed.Seconds(), (1 / elapsed.Seconds()))

	this.HashRate = uint32(float64(count) * (1 / elapsed.Seconds()))

	return
}

func (this *POH) VerifyPOH(cpu_cores int) (err error) {

	start := time.Now()

	// TODO: Revise - Keep one CPU core available for other tasks and scheduling, benchmarks improved.
	if cpu_cores > 4 {
		cpu_cores -= 1
	}

	panicHandler := func(p interface{}) {
		fmt.Printf("Task panicked: %v", p)
	}

	pool := pond.New(cpu_cores, 0, pond.MinWorkers(cpu_cores), pond.PanicHandler(panicHandler))

	// Distribute jobs on each core for the specified sequence
	tasks := uint64(len(this.POH[0].Entry))

	var error_abort bool

	for i := uint64(1); i < tasks; i++ {
		//fmt.Println("Job started => ", i, tasks)
		n := i
		pool.Submit(func() {

			//fmt.Println("Job fork => ", n)

			seqstart := this.POH[0].Entry[n-1].Seq
			seqend := this.POH[0].Entry[n].Seq

			var prevhash []byte

			for a := seqstart; a <= seqend; a++ {

				h := sha256.New()

				// Confirm if the beginning block
				if a == seqstart {
					prevhash = this.POH[0].Entry[n-1].Hash
				} else {
					// Hash the hash
					h.Write([]byte(prevhash))
					prevhash = h.Sum(nil)
				}

				// Confirm the sequence matches our sync state
				if a == seqend {

					// Hash the data block if specified
					if len(this.POH[0].Entry[n].Data) > 0 {
						h.Write(this.POH[0].Entry[n].Data)

						// Verify the signature matches the publickey
						verify, _ := this.Wallet.Verify(this.POH[0].Entry[n].Data, this.POH[0].Entry[n].Signature)

						fmt.Println(verify)

						if verify == false {
							error_abort = true
							log.Printf("POH data signature failed, sequence ID %d - Calculated (%s)", this.POH[0].Entry[n].Seq, this.POH[0].Entry[n].Signature)
						}
					}

					// Read lock
					this.Mu.RLock()
					orig := this.POH[0].Entry[n].Hash
					this.Mu.RUnlock()

					// Compare the proof to the original
					compare := bytes.Compare(h.Sum(nil), orig)

					if compare != 0 {
						error_abort = true
						log.Printf("POH Verification failed, sequence ID %d - Calculated (%s) vs Reference (%s)", this.POH[0].Entry[n].Seq, base64.RawStdEncoding.EncodeToString(h.Sum(nil)), base64.RawStdEncoding.EncodeToString(orig))
						// TODO: Improve stopping existing workers in the queue, revise.
						//pool.Stop()

					}

				}

			}

		})
	}

	// Stop the pool and wait for all submitted tasks to complete
	pool.StopAndWait()

	timer := time.Now()
	elapsed := timer.Sub(start)

	// Calculate the verification hashrate
	lastSeq := this.POH[0].Entry[len(this.POH[0].Entry)-1]
	this.VerifyHashRate = uint32(float64(lastSeq.Seq) * (1 / elapsed.Seconds()))
	this.VerifyHashRatePerCore = this.VerifyHashRate / uint32(cpu_cores)

	if error_abort == true {
		return errors.New("POH validation failed")
	}

	return
}
