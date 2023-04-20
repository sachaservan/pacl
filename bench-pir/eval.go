package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"time"

	paclsk "github.com/sachaservan/pacl/pacl-sk"
	dpf "github.com/sachaservan/vdpf"
)

func main() {

	n := []int{16384, 32768, 65536, 131072, 262144, 524288, 1048576, 2097152, 4194304}
	byteParams := []int{1024, 2048}
	numTrials := 10

	i := 0
	for _, dbSize := range n {
		for _, slotSize := range byteParams {
			slots := make([]*Slot, dbSize)
			for i := 0; i < dbSize; i++ {
				slots[i] = NewRandomSlot(slotSize)
			}
			experiment := &Experiment{
				DBSize:   uint64(dbSize),
				ItemSize: uint64(slotSize),
			}
			experiment.ServerXorProcessingMS = make([]int64, 0)
			experiment.ServerPIRProcessingMS = make([]int64, 0)
			experiment.ServerPIRKeywordProcessingMS = make([]int64, 0)
			experiment.ServerPIRPACLProcessingMS = make([]int64, 0)

			for trial := 0; trial < numTrials; trial++ {
				pirTimeMS, bits := benchmarkPIR(dbSize, slots)
				pirPACLTimeMS := benchmarkPIRPACL(dbSize, slots, bits) // re-use expanded bits to avoid double counting DPF time
				xorTimeMS := benchmarkXor(dbSize, slots, bits)
				pirKeywordTimeMS := benchmarkPIRKeywords(dbSize, slots)
				experiment.ServerXorProcessingMS = append(experiment.ServerXorProcessingMS, xorTimeMS)
				experiment.ServerPIRProcessingMS = append(experiment.ServerPIRProcessingMS, pirTimeMS)
				experiment.ServerPIRKeywordProcessingMS = append(experiment.ServerPIRKeywordProcessingMS, pirKeywordTimeMS)
				experiment.ServerPIRPACLProcessingMS = append(experiment.ServerPIRPACLProcessingMS, pirPACLTimeMS)
				fmt.Printf("Finished trial %v of %v\n", trial, numTrials)
			}

			fmt.Printf("XOR           (%v bytes per item with %v item DB): %v ms\n", slotSize, dbSize, experiment.ServerXorProcessingMS[0])
			fmt.Printf("PIR           (%v bytes per item with %v item DB): %v ms\n", slotSize, dbSize, experiment.ServerPIRProcessingMS[0])
			fmt.Printf("PIR Keyword   (%v bytes per item with %v item DB): %v ms\n", slotSize, dbSize, experiment.ServerPIRKeywordProcessingMS[0])
			fmt.Printf("PIRPACL       (%v bytes per item with %v item DB): %v ms\n", slotSize, dbSize, experiment.ServerPIRPACLProcessingMS[0])

			experimentJSON, err := json.MarshalIndent(experiment, "", " ")
			if err != nil {
				panic(err)
			}
			ioutil.WriteFile("experiment"+fmt.Sprint(i)+".json", experimentJSON, 0644)
			fmt.Println("File saved.")
			i++
		}
	}

}

func benchmarkXor(dbSize int, slots []*Slot, bits []byte) int64 {

	start := time.Now()
	accumulator := NewEmptySlot(len(slots[0].Data))
	for i := 0; i < len(slots); i++ {
		if bits[i]%2 == 1 {
			XorSlots(accumulator, slots[i])
		}
	}

	return time.Since(start).Milliseconds()
}

func benchmarkPIR(dbSize int, slots []*Slot) (int64, []byte) {
	prfKey := dpf.GeneratePRFKey()
	client := dpf.ClientDPFInitialize(prfKey)
	bits := uint(math.Ceil(math.Log2(float64(dbSize))))
	keyA, _ := client.GenDPFKeys(0, bits)
	server := dpf.ServerDPFInitialize(client.PrfKey)

	start := time.Now()

	shares := server.FullDomainEval(keyA)

	accumulator := NewEmptySlot(len(slots[0].Data))
	for i := 0; i < len(slots); i++ {
		if shares[i]%2 == 1 {
			XorSlots(accumulator, slots[i])
		}
	}

	return time.Since(start).Milliseconds(), shares
}

func benchmarkPIRKeywords(dbSize int, slots []*Slot) int64 {
	prfKey := dpf.GeneratePRFKey()
	client := dpf.ClientDPFInitialize(prfKey)

	// PIR-by-keywords requires setting DPF domain to 128 bits for security
	keyA, _ := client.GenDPFKeys(12345, 128)
	server := dpf.ServerDPFInitialize(client.PrfKey)

	// precompute the random values and DPF inputs.
	// Note: DPF evaluation time isn't dependent on input
	x := make([]uint64, dbSize)
	for i := 0; i < dbSize; i++ {
		x[i] = uint64(i) * uint64(i) % uint64(dbSize)
	}

	start := time.Now()

	shares := server.BatchEval(keyA, x)

	accumulator := NewEmptySlot(len(slots[0].Data))
	for i := 0; i < len(slots); i++ {
		if shares[i]%2 == 1 {
			XorSlots(accumulator, slots[i])
		}
	}

	return time.Since(start).Milliseconds()
}

func benchmarkPIRPACL(dbSize int, slots []*Slot, bits []byte) int64 {
	// we can resuse the DPF expansion performed in
	// cloak for PIR so only measure the xor time
	start := time.Now()
	accumulator := NewEmptySlot(len(slots[0].Data))
	for i := 0; i < len(slots); i++ {
		if bits[i]%2 == 1 {
			XorSlots(accumulator, slots[i])
		}
	}
	xortime := time.Since(start).Milliseconds()

	return benchmarkPACL(dbSize) + xortime
}

func benchmarkPACL(dbsize int) int64 {
	// setup parameters
	n := uint(math.Log2(float64(dbsize)))
	kl, key, _ := paclsk.GenerateBenchmarkKeyList(uint64(dbsize), n, paclsk.Equality, 0)

	shares := kl.NewProof(0, key)
	auditB := kl.Audit(shares[1])

	start := time.Now()

	// audit (includes VDPF expansion)
	auditA := kl.Audit(shares[0])
	kl.CheckAudit(auditA, auditB)

	return time.Since(start).Milliseconds()
}
