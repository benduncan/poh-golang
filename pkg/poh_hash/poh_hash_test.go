package poh_hash_test

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/benduncan/poh-golang/pkg/poh_hash"
	"github.com/stretchr/testify/assert"
)

var geneisis_hash string
var count = 10_000_000
var tick_size = 1_000_000

var wallet_path = "../../config/test-wallet.json"

func TestDataVerification(t *testing.T) {

	poh := poh_hash.New(wallet_path)
	data, err := os.ReadFile("../../config/sync-data-validation.json")

	if err != nil {
		t.Error(err)
	}

	var validator []poh_hash.POH_Entry
	err = json.Unmarshal(data, &validator)

	if err != nil {
		t.Error(err)
	}

	poh.POH = append(poh.POH, poh_hash.POH_Epoch{Epoch: 1})
	poh.POH[0].Entry = append(poh.POH[0].Entry, validator...)
	err = poh.VerifyPOH(runtime.NumCPU())

	assert.Nil(t, err)

	// Cause a failure, check hash verification fails
	poh.POH[0].Entry[1].Hash = []byte("9l7h+b1llzUkGfmn+vpSf2btbohK4DISO6KqOmNtJPI=")
	err = poh.VerifyPOH(runtime.NumCPU())
	assert.NotNil(t, err)

}

func TestVerification(t *testing.T) {

	poh := poh_hash.New(wallet_path)
	data, err := os.ReadFile("../../config/sync-validation.json")

	if err != nil {
		t.Error(err)
	}

	var validator []poh_hash.POH_Entry

	err = json.Unmarshal(data, &validator)

	if err != nil {
		t.Error(err)
	}

	poh.POH = append(poh.POH, poh_hash.POH_Epoch{Epoch: 1})
	poh.POH[0].Entry = append(poh.POH[0].Entry, validator...)
	err = poh.VerifyPOH(runtime.NumCPU())

	assert.Nil(t, err)

	// Cause a failure, check hash verification fails
	poh.POH[0].Entry[1].Hash = []byte("9l7h+b1llzUkGfmn+vpSf2btbohK4DISO6KqOmNtJPI=")
	err = poh.VerifyPOH(runtime.NumCPU())
	assert.NotNil(t, err)

}

func TestGenerationVerify(t *testing.T) {

	poh := poh_hash.New(wallet_path)
	poh.QueueSync.State = make([]poh_hash.QueueData, 0)
	go func() {

		// Test periodically adding data to the sync state
		for i := 0; i < 10; i++ {
			time.Sleep(time.Millisecond * (50 * time.Duration(i)))
			str := fmt.Sprintf("Sync state %d", i)
			poh.QueueSync.State = append(poh.QueueSync.State, poh_hash.QueueData{Data: []byte(str)})
		}

	}()

	poh.GeneratePOH(uint64(count))

	err := poh.VerifyPOH(runtime.NumCPU())

	assert.Nil(t, err)

	//	spew.Dump(poh)

	// Create an invalid signature, confirm breaks
	for i := 0; i < len(poh.POH[0].Entry); i++ {

		if len(poh.POH[0].Entry[i].Data) > 0 {
			poh.POH[0].Entry[i].Signature = []byte("invalid")
			break
		}

	}

	err = poh.VerifyPOH(runtime.NumCPU())

	assert.NotNil(t, err)

}

func BenchmarkGeneratePOH_10000(b *testing.B) {

	for n := 0; n < b.N; n++ {
		poh := poh_hash.New("")
		poh.GeneratePOH(10_000)
	}

}

func BenchmarkGeneratePOH_100000(b *testing.B) {

	for n := 0; n < b.N; n++ {
		poh := poh_hash.New("")
		poh.GeneratePOH(100_000)
	}

}

func BenchmarkGeneratePOH_1000000(b *testing.B) {

	for n := 0; n < b.N; n++ {
		poh := poh_hash.New("")
		poh.GeneratePOH(1_000_000)
	}

}

func BenchmarkGenerateDataPOH_100000(b *testing.B) {

	for n := 0; n < b.N; n++ {
		poh := poh_hash.New("")
		poh.GeneratePOH(100_000)

		go func() {

			// Test periodically adding data to the sync state
			for i := 0; i < 10_000; i++ {
				time.Sleep(time.Microsecond * (1 * time.Duration(i)))
				str := fmt.Sprintf("Sync state %d", i)
				poh.QueueSync.State = append(poh.QueueSync.State, poh_hash.QueueData{Data: []byte(str)})
			}

		}()

	}

}

func BenchmarkGenerateDataPOH_1000000(b *testing.B) {

	for n := 0; n < b.N; n++ {
		poh := poh_hash.New("")
		poh.GeneratePOH(1_000_000)

		go func() {

			// Test periodically adding data to the sync state
			for i := 0; i < 100_000; i++ {
				time.Sleep(time.Microsecond * (1 * time.Duration(i)))
				str := fmt.Sprintf("Sync state %d", i)
				poh.QueueSync.State = append(poh.QueueSync.State, poh_hash.QueueData{Data: []byte(str)})
			}

		}()

	}

}

func BenchmarkVerifyPOH_AllCores(b *testing.B) {

	cpu_cores := runtime.NumCPU()

	poh := poh_hash.New("")
	poh.GeneratePOH(10_000)

	for n := 0; n < b.N; n++ {

		poh.VerifyPOH(cpu_cores)

	}

}

func BenchmarkVerifyPOH_AllCoresMinusOne(b *testing.B) {

	cpu_cores := runtime.NumCPU()

	poh := poh_hash.New("")
	poh.GeneratePOH(10_000)

	for n := 0; n < b.N; n++ {

		poh.VerifyPOH(cpu_cores - 1)

	}

}
func BenchmarkVerifyPOH_AllCores_Double(b *testing.B) {

	cpu_cores := runtime.NumCPU()

	poh := poh_hash.New("")
	poh.GeneratePOH(10_000)

	for n := 0; n < b.N; n++ {

		poh.VerifyPOH(cpu_cores * 2)

	}

}

func BenchmarkVerifyPOH_QuadCore(b *testing.B) {

	poh := poh_hash.New("")
	poh.GeneratePOH(10_000)

	for n := 0; n < b.N; n++ {

		poh.VerifyPOH(4)

	}

}

func BenchmarkVerifyPOH_OctCore(b *testing.B) {

	poh := poh_hash.New("")
	poh.GeneratePOH(10_000)

	for n := 0; n < b.N; n++ {
		poh.VerifyPOH(8)
	}

}
