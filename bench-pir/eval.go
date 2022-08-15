package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"time"

	"github.com/sachaservan/pacl/algebra"
	paclsposs "github.com/sachaservan/pacl/pacl-sposs"
	"github.com/sachaservan/pacl/sposs"
	dpf "github.com/sachaservan/vdpf"
)

func main() {

	n := []int{16384, 32768, 65536, 131072, 262144, 524288, 1048576, 2097152, 4194304}
	byteParams := []int{256, 512, 1024}

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
			experiment.ServerPIRProcessingMS = make([]int64, 0)
			experiment.ServerPIRPACLProcessingMS = make([]int64, 0)

			for trial := 0; trial < 10; trial++ {
				pirTimeMS, bits := benchmarkPIR(dbSize, slots)
				pirPACLTimeMS := benchmarkPIRPACL(dbSize, slots, bits)
				experiment.ServerPIRProcessingMS = append(experiment.ServerPIRProcessingMS, pirTimeMS)
				experiment.ServerPIRPACLProcessingMS = append(experiment.ServerPIRPACLProcessingMS, pirPACLTimeMS)
				fmt.Printf("Finished trial %v of %v\n", trial, 10)
			}

			fmt.Printf("PIR       (%v bytes per item with %v item DB): %v\n", slotSize, dbSize, experiment.ServerPIRProcessingMS[0])
			fmt.Printf("PIRPACL   (%v bytes per item with %v item DB): %v\n", slotSize, dbSize, experiment.ServerPIRPACLProcessingMS[0])

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
	domain := uint(math.Log2(float64(dbsize) + 1))
	group := paclsposs.DefaultGroup()
	kl, key, _ := paclsposs.GenerateBenchmarkKeyList(uint64(dbsize), domain, group)

	shares := kl.NewProof(0, key)

	start := time.Now()

	// prepare the audit
	kl.Audit(shares[0])

	// simulate verification of the sposs proof
	spossTime := benchmarkSPoSSProofMS(group)

	return time.Since(start).Milliseconds() + spossTime
}

func benchmarkSPoSSProofMS(group *algebra.Group) int64 {
	start := time.Now()

	pp := sposs.NewPublicParams(group)
	gX := pp.Group.Field.MulIdentity()
	additiveShareA, _ := pp.LinearShares(gX)
	proofA, _ := pp.GenProof(pp.ExpField.AddIdentity())
	auditShareA := pp.Audit(additiveShareA, proofA)
	pp.CheckAudit(auditShareA, auditShareA)

	return time.Since(start).Milliseconds()
}
