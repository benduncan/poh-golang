package poh_hash_test

import (
	"runtime"
	"testing"

	"github.com/benduncan/poh-golang/pkg/poh_hash"
)

func BenchmarkGeneratePOH_10000(b *testing.B) {

	for n := 0; n < b.N; n++ {
		poh := poh_hash.New()
		poh.GeneratePOH(10000)
	}

}

func BenchmarkGeneratePOH_100000(b *testing.B) {

	for n := 0; n < b.N; n++ {
		poh := poh_hash.New()
		poh.GeneratePOH(100_000)
	}

}

func BenchmarkGeneratePOH_1000000(b *testing.B) {

	for n := 0; n < b.N; n++ {
		poh := poh_hash.New()
		poh.GeneratePOH(1_000_000)
	}

}

func BenchmarkVerifyPOH_AllCores(b *testing.B) {

	cpu_cores := runtime.NumCPU()

	poh := poh_hash.New()
	poh.GeneratePOH(10_000)

	for n := 0; n < b.N; n++ {

		poh.VerifyPOH(cpu_cores)

	}

}

func BenchmarkVerifyPOH_AllCoresMinusOne(b *testing.B) {

	cpu_cores := runtime.NumCPU()

	poh := poh_hash.New()
	poh.GeneratePOH(10_000)

	for n := 0; n < b.N; n++ {

		poh.VerifyPOH(cpu_cores - 1)

	}

}
func BenchmarkVerifyPOH_AllCores_Double(b *testing.B) {

	cpu_cores := runtime.NumCPU()

	poh := poh_hash.New()
	poh.GeneratePOH(10_000)

	for n := 0; n < b.N; n++ {

		poh.VerifyPOH(cpu_cores * 2)

	}

}

func BenchmarkVerifyPOH_QuadCore(b *testing.B) {

	poh := poh_hash.New()
	poh.GeneratePOH(10_000)

	for n := 0; n < b.N; n++ {

		poh.VerifyPOH(4)

	}

}

func BenchmarkVerifyPOH_OctCore(b *testing.B) {

	poh := poh_hash.New()
	poh.GeneratePOH(10_000)

	for n := 0; n < b.N; n++ {
		poh.VerifyPOH(8)
	}

}
