package main

import (
	"crypto/elliptic"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	paclpk "github.com/sachaservan/pacl/pacl-pk"
	paclsk "github.com/sachaservan/pacl/pacl-sk"
	paclsposs "github.com/sachaservan/pacl/pacl-sposs"
	dpf "github.com/sachaservan/vdpf"
)

const (
	FSSType     int = iota
	FSSEquality     = 0
	FSSRange        = 1
)

type BaselineKey struct {
	Verifiable bool
	HashKeys   [2]dpf.HashKey
	PRFKey     dpf.PrfKey
	DPFKey     *dpf.DPFKey
	Indices    []uint64
}

func main() {

	numTrials := 1000

	// amortize the proof verification across numEval
	numEvalKeys := []uint64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}
	numEvalSubkeys := []uint64{1, 10, 20}
	domain := []uint{32}

	i := 0
	for _, fssDomain := range domain {
		for _, numKeys := range numEvalKeys {
			for _, numSubkeys := range numEvalSubkeys {

				amortization := int64(numKeys)

				//////////////////////////////////
				// generate the baseline DPF and VDPF keys
				//////////////////////////////////
				baselineIndices := make([]uint64, numKeys)
				prfKey := dpf.GeneratePRFKey()
				hashKeys := dpf.GenerateVDPFHashKeys()

				// init DPF
				pf := dpf.ClientDPFInitialize(prfKey)
				keyDPF, _ := pf.GenDPFKeys(1, fssDomain)

				// init VDPF
				pfver := dpf.ClientVDPFInitialize(prfKey, hashKeys)
				keyVDPF, _ := pfver.GenDPFKeys(1, fssDomain)

				for i := 0; i < int(numKeys); i++ {
					baselineIndices[i] = rand.Uint64()
				}

				baselineDPFKey := &BaselineKey{
					Verifiable: false,
					PRFKey:     prfKey,
					DPFKey:     keyDPF,
					Indices:    baselineIndices,
				}

				baselineVDPFKey := &BaselineKey{
					Verifiable: true,
					HashKeys:   hashKeys,
					PRFKey:     prfKey,
					DPFKey:     keyVDPF,
					Indices:    baselineIndices,
				}
				//////////////////////////////////

				//////////////////////////////////
				// setup the keylists for this set of parameters
				//////////////////////////////////
				klpk, xpk, idxpk := paclpk.GenerateBenchmarkKeyList(numKeys, fssDomain, elliptic.P256(), paclpk.Inclusion, numSubkeys)
				klsk, xsk, idxsk := paclsk.GenerateBenchmarkKeyList(numKeys, fssDomain, paclsk.Inclusion, numSubkeys)
				klsposs, xsposs, idxsposs := paclsposs.GenerateBenchmarkKeyList(
					numKeys, fssDomain, paclsposs.DefaultGroup(), paclsposs.Inclusion, numSubkeys)

				// generate the PACL proofs
				sharesPk := klpk.NewProof(idxpk, xpk)
				sharesSk := klsk.NewProof(idxsk, xsk)
				sharesSposs := klsposs.NewProof(idxsposs, xsposs)
				//////////////////////////////////

				// initialize the experiment for this set of parameters
				experiment := &Experiment{
					FSSDomain:  uint64(fssDomain),
					NumKeys:    numKeys,
					NumSubkeys: numSubkeys,
				}

				// WARMUP: do a trial run as a warmup
				for trial := 0; trial < numTrials; trial++ {
					benchmarkBaselineFSS(baselineDPFKey, fssDomain, FSSRange)
					benchmarkBaselineFSS(baselineVDPFKey, fssDomain, FSSRange)
					benchmarkPACLPublicKeyFSS(baselineDPFKey, klpk, sharesPk[0], fssDomain, FSSRange)
					benchmarkPACLSymmetricKeyFSS(baselineDPFKey, klsk, sharesSk[0], fssDomain, FSSRange)
				}

				// measure group exponentiation time
				_, x := klsposs.Group.RandomElement()
				xF := klsposs.ProofPP.ExpField.NewElement(x)
				timeExp := time.Now()
				gX := klsposs.Group.NewElement(x)

				for trial := 0; trial < numTrials; trial++ {
					gX = klsposs.Group.Exp(gX, xF) // group exponentiation
				}

				experiment.GroupExponentiationNS = uint64(time.Since(timeExp)) / uint64(amortization)

				// (V)DPF evaluation baselines
				for trial := 0; trial < numTrials; trial++ {
					// DPF
					timeEq := benchmarkBaselineFSS(baselineDPFKey, fssDomain, FSSEquality)
					timeEq /= amortization
					experiment.EqualityBaselineProcessingNS = append(experiment.EqualityBaselineProcessingNS, timeEq)

					// VDPF
					timeEq = benchmarkBaselineFSS(baselineVDPFKey, fssDomain, FSSEquality)
					timeEq /= amortization
					experiment.EqualityBaselineVerProcessingNS = append(experiment.EqualityBaselineVerProcessingNS, timeEq)

					// DMPF range
					timeRange := benchmarkBaselineFSS(baselineDPFKey, fssDomain, FSSRange)
					timeRange /= amortization
					experiment.RangeBaselineProcessingNS = append(experiment.RangeBaselineProcessingNS, timeRange)

					// VDMPF range
					timeRange = benchmarkBaselineFSS(baselineVDPFKey, fssDomain, FSSRange)
					timeRange /= amortization
					experiment.RangeBaselineVerProcessingNS = append(experiment.RangeBaselineVerProcessingNS, timeRange)
				}

				// DPF PACL (public key)
				for trial := 0; trial < numTrials; trial++ {
					// DPF with public key PACL
					timeEq := benchmarkPACLPublicKeyFSS(baselineDPFKey, klpk, sharesPk[0], fssDomain, FSSEquality)
					timeEq /= amortization
					experiment.EqualityDPFPACLProcessingNS = append(experiment.EqualityDPFPACLProcessingNS, timeEq)

					// DMPF range with public key PACL
					timeRange := benchmarkPACLPublicKeyFSS(baselineDPFKey, klpk, sharesPk[0], fssDomain, FSSRange)
					timeRange /= amortization
					experiment.RangeDPFPACLProcessingNS = append(experiment.RangeDPFPACLProcessingNS, timeRange)
				}

				// VDPF PACL (public key sposs)
				for trial := 0; trial < numTrials; trial++ {
					// equality
					timeEq := benchmarkPACLPublicKeyVFSS(baselineVDPFKey, klsposs, sharesSposs[0], fssDomain, FSSEquality)
					timeEq /= amortization
					experiment.EqualityVDPFPACLProcessingNS = append(experiment.EqualityVDPFPACLProcessingNS, timeEq)

					// range
					timeRange := benchmarkPACLPublicKeyVFSS(baselineVDPFKey, klsposs, sharesSposs[0], fssDomain, FSSRange)
					timeRange /= amortization
					experiment.RangeVDPFPACLProcessingNS = append(experiment.RangeVDPFPACLProcessingNS, timeRange)
				}

				// DPF PACL (symmetric key)
				for trial := 0; trial < numTrials; trial++ {
					// equality
					timeEq := benchmarkPACLSymmetricKeyFSS(baselineDPFKey, klsk, sharesSk[0], fssDomain, FSSEquality)
					timeEq /= amortization
					experiment.EqualityDPFSKPACLProcessingNS = append(experiment.EqualityDPFSKPACLProcessingNS, timeEq)

					// range
					timeRange := benchmarkPACLSymmetricKeyFSS(baselineDPFKey, klsk, sharesSk[0], fssDomain, FSSRange)
					timeRange /= amortization
					experiment.RangeDPFSKPACLProcessingNS = append(experiment.RangeDPFSKPACLProcessingNS, timeRange)
				}

				// VDPF PACL (symmetric key)
				for trial := 0; trial < numTrials; trial++ {
					// equality
					timeEq := benchmarkPACLSymmetricKeyVFSS(baselineVDPFKey, klsposs, sharesSposs[0], fssDomain, FSSEquality)
					timeEq /= amortization
					experiment.EqualityVDPFSKPACLProcessingNS = append(experiment.EqualityVDPFSKPACLProcessingNS, timeEq)

					// range
					timeRange := benchmarkPACLSymmetricKeyVFSS(baselineVDPFKey, klsposs, sharesSposs[0], fssDomain, FSSRange)
					timeRange /= amortization
					experiment.RangeVDPFSKPACLProcessingNS = append(experiment.RangeVDPFSKPACLProcessingNS, timeRange)
				}

				fmt.Println("---------------------------------")
				fmt.Printf("FSS domain:     %v\n", fssDomain)
				fmt.Printf("Num keys:       %v\n", numKeys)
				fmt.Printf("Num subkeys:    %v\n", numSubkeys)
				fmt.Printf("Exponentiation: %v\n", experiment.GroupExponentiationNS))
				fmt.Println("---------------------------------")
				fmt.Printf("DPF (x = a)              (size %v): %v\n", fssDomain, avg(experiment.EqualityBaselineProcessingNS))
				fmt.Printf("DPF SK-PACL (x = a)      (size %v): %v\n", fssDomain, avg(experiment.EqualityDPFSKPACLProcessingNS))
				fmt.Printf("DPF PACL (x = a)         (size %v): %v\n", fssDomain, avg(experiment.EqualityDPFPACLProcessingNS))
				fmt.Printf("VDPF (x = a)             (size %v): %v\n", fssDomain, avg(experiment.EqualityBaselineVerProcessingNS))
				fmt.Printf("VDPF SK PACL (x = a)     (size %v): %v\n", fssDomain, avg(experiment.EqualityVDPFSKPACLProcessingNS))
				fmt.Printf("VDPF PACL (x = a)        (size %v): %v\n", fssDomain, avg(experiment.EqualityVDPFPACLProcessingNS))
				fmt.Println("---------------------------------")
				fmt.Printf("DPF (a < x < b)          (size %v): %v\n", fssDomain, avg(experiment.RangeBaselineProcessingNS))
				fmt.Printf("DPF SK-PACL (a < x < b)  (size %v): %v\n", fssDomain, avg(experiment.RangeDPFSKPACLProcessingNS))
				fmt.Printf("DPF PACL (a < x < b)     (size %v): %v\n", fssDomain, avg(experiment.RangeDPFPACLProcessingNS))
				fmt.Printf("VDPF (a < x < b)         (size %v): %v\n", fssDomain, avg(experiment.RangeBaselineVerProcessingNS))
				fmt.Printf("VDPF SK PACL (a < x < b) (size %v): %v\n", fssDomain, avg(experiment.RangeVDPFSKPACLProcessingNS))
				fmt.Printf("VDPF PACL (a < x < b)    (size %v): %v\n", fssDomain, avg(experiment.RangeVDPFPACLProcessingNS))
				fmt.Println("---------------------------------")

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
}

func avg(arr []int64) float64 {
	sum := int64(0)
	for i := 0; i < len(arr); i++ {
		sum += arr[i]
	}

	return float64(sum) / float64(len(arr))
}

func benchmarkBaselineFSS(
	key *BaselineKey,
	fssDomain uint,
	fssType int) int64 {

	totalTime := int64(0)

	if fssType == FSSEquality {
		// equality is just DPF
		key.DPFKey = randomizeDPFKey(key.DPFKey)
		start := time.Now()
		expandBaselineDPF(key)
		totalTime += time.Since(start).Microseconds()

	} else if fssType == FSSRange {
		// (blackbox) inequality is 2*logn invocations of a DPF
		for i := uint(0); i < 2*fssDomain; i++ {
			key.DPFKey = randomizeDPFKey(key.DPFKey)
			start := time.Now()
			expandBaselineDPF(key)
			totalTime += time.Since(start).Microseconds()
		}
	}

	return totalTime
}

func benchmarkPACLPublicKeyFSS(
	key *BaselineKey,
	kl *paclpk.KeyList,
	share *paclpk.ProofShare,
	fssDomain uint,
	fssType int) int64 {

	totalTime := int64(0)

	if fssType == FSSRange {
		// (blackbox) inequality is 2*logn invocations of a DPF
		// 2*fssDomain-1 because the last DPF evaluation is implicit in Audit
		for i := uint(0); i < 2*fssDomain-1; i++ {
			key.DPFKey = randomizeDPFKey(key.DPFKey)
			start := time.Now()
			expandBaselineDPF(key)
			totalTime += time.Since(start).Microseconds()
		}
	}

	share.DPFKey = randomizeDPFKey(share.DPFKey)
	start := time.Now()
	kl.Audit(share)
	totalTime += time.Since(start).Microseconds()

	return totalTime
}

func benchmarkPACLSymmetricKeyFSS(
	key *BaselineKey,
	kl *paclsk.KeyList,
	share *paclsk.ProofShare,
	fssDomain uint,
	fssType int) int64 {

	totalTime := int64(0)

	if fssType == FSSRange {
		// (blackbox) inequality is 2*logn invocations of a DPF
		// 2*fssDomain-1 because the last DPF evaluation is implicit in Audit
		for i := uint(0); i < 2*fssDomain-1; i++ {
			key.DPFKey = randomizeDPFKey(key.DPFKey)
			start := time.Now()
			expandBaselineDPF(key)
			totalTime += time.Since(start).Microseconds()
		}
	}

	share.DPFKey = randomizeDPFKey(share.DPFKey)
	start := time.Now()
	kl.Audit(share)
	totalTime += time.Since(start).Microseconds()

	return totalTime
}

func benchmarkPACLSymmetricKeyVFSS(
	key *BaselineKey,
	kl *paclsposs.KeyList,
	share *paclsposs.ProofShare,
	fssDomain uint,
	fssType int) int64 {

	totalTime := int64(0)

	// TODO: this benchmark has sloppy edge cases
	// because we want to run the VDPF but then not run SPoSS so what happens
	// is that we first run the Expand (which expands to N*l values)
	// and then compute the xors with the keys.
	// This could use some cleaned up

	if fssType == FSSRange {
		// (blackbox) inequality is 2*logn invocations of a DPF
		// 2*fssDomain-1 because the last DPF evaluation is implicit in Audit
		for i := uint(0); i < 2*fssDomain-1; i++ {
			key.DPFKey = randomizeDPFKey(key.DPFKey)
			start := time.Now()
			expandBaselineDPF(key)
			totalTime += time.Since(start).Microseconds()
		}
	}

	share.DPFKey = randomizeDPFKey(share.DPFKey)
	start := time.Now()
	bits, _ := kl.ExpandVDPF(share)
	totalTime += time.Since(start).Microseconds()

	// make a bunch of random symmetric keys
	// bits is suppposed to si
	slots := make([]*paclsk.Slot, kl.NumKeys)
	for i := 0; i < len(slots); i++ {
		slots[i] = paclsk.NewRandomSlot(16) // symmetric key is 16 bytes
	}

	start = time.Now()
	accumulator := paclsk.NewEmptySlot(16)
	for i := 0; i < len(slots); i++ {
		if bits[i]%2 == 1 {
			paclsk.XorSlots(accumulator, slots[i])
		}
	}
	totalTime += time.Since(start).Microseconds()

	return totalTime
}

func benchmarkPACLPublicKeyVFSS(
	key *BaselineKey,
	kl *paclsposs.KeyList,
	share *paclsposs.ProofShare,
	fssDomain uint,
	fssType int) int64 {

	totalTime := int64(0)

	if fssType == FSSRange {
		// (blackbox) inequality is 2*logn invocations of a DPF
		// 2*fssDomain-1 because the last DPF evaluation is implicit in Audit
		for i := uint(0); i < 2*fssDomain-1; i++ {
			key.DPFKey = randomizeDPFKey(key.DPFKey)
			start := time.Now()
			expandBaselineDPF(key)
			totalTime += time.Since(start).Microseconds()
		}
	}

	share.DPFKey = randomizeDPFKey(share.DPFKey)
	start := time.Now()
	kl.Audit(share)
	totalTime += time.Since(start).Microseconds()

	return totalTime
}

func expandBaselineDPF(key *BaselineKey) {
	if key.Verifiable {
		pf := dpf.ServerVDPFInitialize(key.PRFKey, key.HashKeys)
		pf.BatchVerEval(key.DPFKey, key.Indices)
	} else {
		pf := dpf.ServerDPFInitialize(key.PRFKey)
		pf.BatchEval(key.DPFKey, key.Indices)
	}
}

func randomizeDPFKey(dpfKey *dpf.DPFKey) *dpf.DPFKey {
	// super hacky way to generate a random gibberish DPF key
	// but it's sufficient for accurate benchmarks
	// see ../vdpf/wrapper.go
	keySize := 18*dpfKey.RangeSize + 18 + 16 + 16*4
	r := make([]byte, keySize)
	_, _ = rand.Read(r)
	for i := 0; i < int(keySize); i++ {
		dpfKey.Bytes[i] = r[i]
	}

	return dpfKey
}
