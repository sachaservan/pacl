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
	FSSType       int = iota
	FSSEquality       = 0
	FSSInequality     = 1
	FSSRange          = 2
	FSSDTree          = 3
)

func main() {

	numTrials := 10

	// amortize the proof verification across numEval
	numEvalKeys := []uint64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}
	numEvalSubkeys := []uint64{1, 10, 20}
	domain := []uint{32}

	i := 0
	for _, fssDomain := range domain {
		for _, numKeys := range numEvalKeys {
			for _, numSubkeys := range numEvalSubkeys {

				klpk, xpk, _ := paclpk.GenerateBenchmarkKeyList(numKeys, fssDomain, elliptic.P256(), paclpk.Inclusion, numSubkeys)
				klsk, xsk, _ := paclsk.GenerateBenchmarkKeyList(numKeys, fssDomain, paclsk.Inclusion, numSubkeys)
				klsposs, xsposs, _ := paclsposs.GenerateBenchmarkKeyList(
					numKeys, fssDomain, paclsposs.DefaultGroup(), paclsposs.Inclusion, numSubkeys)

				sharesPk := klpk.NewProof(0, xpk)
				sharesSk := klsk.NewProof(0, xsk)
				sharesSposs := klsposs.NewProof(0, xsposs)

				experiment := &Experiment{
					FSSDomain: uint64(fssDomain),
					NumKeys:   numKeys,
				}

				for trial := 0; trial < numTrials; trial++ {
					// WARMUP
					benchmarkBaselineFSS(klpk, sharesPk[0], fssDomain, FSSEquality)
					benchmarkBaselineVFSS(klsposs, sharesSposs[0], fssDomain, FSSEquality)
					benchmarkPACLSymmetricKeyFSS(klsk, sharesSk[0], fssDomain, FSSEquality)
				}

				// (V)DPF evaluation baseline
				for trial := 0; trial < numTrials; trial++ {
					// equality
					timeEq := benchmarkBaselineFSS(klpk, sharesPk[0], fssDomain, FSSEquality)
					experiment.EqualityBaselineProcessingNS = append(experiment.EqualityBaselineProcessingNS, timeEq)

					timeEq = benchmarkBaselineVFSS(klsposs, sharesSposs[0], fssDomain, FSSEquality)
					experiment.EqualityBaselineVerProcessingNS = append(experiment.EqualityBaselineVerProcessingNS, timeEq)

					// inequality
					timeIneq := benchmarkBaselineFSS(klpk, sharesPk[0], fssDomain, FSSInequality)
					experiment.InequalityBaselineProcessingNS = append(experiment.InequalityBaselineProcessingNS, timeIneq)

					timeIneq = benchmarkBaselineVFSS(klsposs, sharesSposs[0], fssDomain, FSSInequality)
					experiment.InequalityBaselineVerProcessingNS = append(experiment.InequalityBaselineVerProcessingNS, timeIneq)

					// range
					timeRange := benchmarkBaselineFSS(klpk, sharesPk[0], fssDomain, FSSRange)
					experiment.RangeBaselineProcessingNS = append(experiment.RangeBaselineProcessingNS, timeRange)

					timeRange = benchmarkBaselineVFSS(klsposs, sharesSposs[0], fssDomain, FSSRange)
					experiment.RangeBaselineVerProcessingNS = append(experiment.RangeBaselineVerProcessingNS, timeRange)

				}

				// DPF PACL (public key)
				for trial := 0; trial < numTrials; trial++ {
					// equality
					timeEq := benchmarkPACLPublicKeyFSS(klpk, sharesPk[0], fssDomain, FSSEquality)
					experiment.EqualityDPFPACLProcessingNS = append(experiment.EqualityDPFPACLProcessingNS, timeEq)

					// inequality
					timeIneq := benchmarkPACLPublicKeyFSS(klpk, sharesPk[0], fssDomain, FSSInequality)
					experiment.InequalityDPFPACLProcessingNS = append(experiment.InequalityDPFPACLProcessingNS, timeIneq)

					// range
					timeRange := benchmarkPACLPublicKeyFSS(klpk, sharesPk[0], fssDomain, FSSRange)
					experiment.RangeDPFPACLProcessingNS = append(experiment.RangeDPFPACLProcessingNS, timeRange)
				}

				// VDPF PACL (public key sposs)
				for trial := 0; trial < numTrials; trial++ {
					// equality
					timeEq := benchmarkPACLPublicKeyVFSS(klsposs, sharesSposs[0], fssDomain, FSSEquality)
					experiment.EqualityVDPFPACLProcessingNS = append(experiment.EqualityVDPFPACLProcessingNS, timeEq)

					// inequality
					timeIneq := benchmarkPACLPublicKeyVFSS(klsposs, sharesSposs[0], fssDomain, FSSInequality)
					experiment.InequalityVDPFPACLProcessingNS = append(experiment.InequalityVDPFPACLProcessingNS, timeIneq)

					// range
					timeRange := benchmarkPACLPublicKeyVFSS(klsposs, sharesSposs[0], fssDomain, FSSRange)
					experiment.RangeVDPFPACLProcessingNS = append(experiment.RangeVDPFPACLProcessingNS, timeRange)
				}

				// DPF PACL (symmetric key)
				for trial := 0; trial < numTrials; trial++ {
					// equality
					timeEq := benchmarkPACLSymmetricKeyFSS(klsk, sharesSk[0], fssDomain, FSSEquality)
					experiment.EqualityDPFSKPACLProcessingNS = append(experiment.EqualityDPFSKPACLProcessingNS, timeEq)

					// inequality
					timeIneq := benchmarkPACLSymmetricKeyFSS(klsk, sharesSk[0], fssDomain, FSSInequality)
					experiment.InequalityDPFSKPACLProcessingNS = append(experiment.InequalityDPFSKPACLProcessingNS, timeIneq)

					// range
					timeRange := benchmarkPACLSymmetricKeyFSS(klsk, sharesSk[0], fssDomain, FSSRange)
					experiment.RangeDPFSKPACLProcessingNS = append(experiment.RangeDPFSKPACLProcessingNS, timeRange)
				}

				// VDPF PACL (symmetric key)
				for trial := 0; trial < numTrials; trial++ {
					// equality
					timeEq := benchmarkPACLSymmetricKeyVFSS(klsposs, sharesSposs[0], fssDomain, FSSEquality)
					experiment.EqualityVDPFSKPACLProcessingNS = append(experiment.EqualityVDPFSKPACLProcessingNS, timeEq)

					// inequality
					timeIneq := benchmarkPACLSymmetricKeyVFSS(klsposs, sharesSposs[0], fssDomain, FSSInequality)
					experiment.InequalityVDPFSKPACLProcessingNS = append(experiment.InequalityVDPFSKPACLProcessingNS, timeIneq)

					// range
					timeRange := benchmarkPACLSymmetricKeyVFSS(klsposs, sharesSposs[0], fssDomain, FSSRange)
					experiment.RangeVDPFSKPACLProcessingNS = append(experiment.RangeVDPFSKPACLProcessingNS, timeRange)
				}

				fmt.Println("---------------------------------")
				fmt.Printf("FSS domain:  %v\n", fssDomain)
				fmt.Printf("Num keys:    %v\n", numKeys)
				fmt.Printf("Num subkeys: %v\n", numSubkeys)
				fmt.Println("---------------------------------")
				fmt.Printf("DPF (x = a)              (size %v): %v\n", fssDomain, avg(experiment.EqualityBaselineProcessingNS))
				fmt.Printf("DPF SK-PACL (x = a)      (size %v): %v\n", fssDomain, avg(experiment.EqualityDPFSKPACLProcessingNS))
				fmt.Printf("DPF PACL (x = a)         (size %v): %v\n", fssDomain, avg(experiment.EqualityDPFPACLProcessingNS))
				fmt.Printf("VDPF (x = a)             (size %v): %v\n", fssDomain, avg(experiment.EqualityBaselineVerProcessingNS))
				fmt.Printf("VDPF SK PACL (x = a)     (size %v): %v\n", fssDomain, avg(experiment.EqualityVDPFSKPACLProcessingNS))
				fmt.Printf("VDPF PACL (x = a)        (size %v): %v\n", fssDomain, avg(experiment.EqualityVDPFPACLProcessingNS))
				fmt.Println("---------------------------------")
				fmt.Printf("DPF (x < a)              (size %v): %v\n", fssDomain, avg(experiment.InequalityBaselineProcessingNS))
				fmt.Printf("DPF SK-PACL (x < a)      (size %v): %v\n", fssDomain, avg(experiment.InequalityDPFSKPACLProcessingNS))
				fmt.Printf("DPF PACL (x < a)         (size %v): %v\n", fssDomain, avg(experiment.InequalityDPFPACLProcessingNS))
				fmt.Printf("VDPF (x < a)             (size %v): %v\n", fssDomain, avg(experiment.InequalityBaselineVerProcessingNS))
				fmt.Printf("VDPF SK PACL (x < a)     (size %v): %v\n", fssDomain, avg(experiment.InequalityVDPFSKPACLProcessingNS))
				fmt.Printf("VDPF PACL (x < a)        (size %v): %v\n", fssDomain, avg(experiment.InequalityVDPFPACLProcessingNS))
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
	kl *paclpk.KeyList,
	share *paclpk.ProofShare,
	fssDomain uint,
	fssType int) int64 {

	totalTime := int64(0)

	if fssType == FSSEquality {
		// equality is just DPF
		share.DPFKey = randomizeDPFKey(share.DPFKey)
		start := time.Now()
		kl.ExpandDPF(share)
		totalTime += time.Since(start).Microseconds()
	} else if fssType == FSSInequality {
		// (blackbox) inequality is logn invocations of a DPF
		for i := uint(0); i < fssDomain; i++ {
			share.DPFKey = randomizeDPFKey(share.DPFKey)
			start := time.Now()
			kl.ExpandDPF(share)
			totalTime += time.Since(start).Microseconds()
		}
	} else if fssType == FSSRange {
		// (blackbox) inequality is 2*logn invocations of a DPF
		for i := uint(0); i < 2*fssDomain; i++ {
			share.DPFKey = randomizeDPFKey(share.DPFKey)
			start := time.Now()
			kl.ExpandDPF(share)
			totalTime += time.Since(start).Microseconds()
		}
	} else if fssType == FSSDTree {
		// for each level (node), we evaluate one DPF
		// over a 32 bit domain.
		for i := uint(0); i < 2*fssDomain; i++ {
			//TODO:
		}
	}

	return totalTime
}

func benchmarkBaselineVFSS(
	kl *paclsposs.KeyList,
	share *paclsposs.ProofShare,
	fssDomain uint,
	fssType int) int64 {

	totalTime := int64(0)

	if fssType == FSSEquality {
		// equality is just DPF
		share.DPFKey = randomizeDPFKey(share.DPFKey)
		start := time.Now()
		kl.ExpandVDPF(share)
		totalTime += time.Since(start).Microseconds()
	} else if fssType == FSSInequality {
		// (blackbox) inequality is logn invocations of a DPF
		for i := uint(0); i < fssDomain; i++ {
			share.DPFKey = randomizeDPFKey(share.DPFKey)
			start := time.Now()
			kl.ExpandVDPF(share)
			totalTime += time.Since(start).Microseconds()
		}
	} else if fssType == FSSRange {
		// (blackbox) inequality is 2*logn invocations of a DPF
		for i := uint(0); i < 2*fssDomain; i++ {
			share.DPFKey = randomizeDPFKey(share.DPFKey)
			start := time.Now()
			kl.ExpandVDPF(share)
			totalTime += time.Since(start).Microseconds()
		}
	}

	return totalTime
}

func benchmarkPACLPublicKeyFSS(
	kl *paclpk.KeyList,
	share *paclpk.ProofShare,
	fssDomain uint,
	fssType int) int64 {

	totalTime := int64(0)

	if fssType == FSSEquality {
		// equality is just DPF
		share.DPFKey = randomizeDPFKey(share.DPFKey)

	} else if fssType == FSSInequality {

		// (blackbox) inequality is logn invocations of a DPF
		for i := uint(0); i < fssDomain; i++ {
			share.DPFKey = randomizeDPFKey(share.DPFKey)
			start := time.Now()
			kl.ExpandDPF(share)
			totalTime += time.Since(start).Microseconds()
		}

	} else if fssType == FSSRange {
		// (blackbox) inequality is 2*logn invocations of a DPF
		for i := uint(0); i < 2*fssDomain; i++ {
			share.DPFKey = randomizeDPFKey(share.DPFKey)
			start := time.Now()
			kl.ExpandDPF(share)
			totalTime += time.Since(start).Microseconds()
		}
	}

	start := time.Now()
	kl.Audit(share)
	totalTime += time.Since(start).Microseconds()

	return totalTime
}

func benchmarkPACLSymmetricKeyFSS(
	kl *paclsk.KeyList,
	share *paclsk.ProofShare,
	fssDomain uint,
	fssType int) int64 {

	totalTime := int64(0)

	if fssType == FSSEquality {
		// equality is just DPF
		share.DPFKey = randomizeDPFKey(share.DPFKey)

	} else if fssType == FSSInequality {

		// (blackbox) inequality is logn invocations of a DPF
		for i := uint(0); i < fssDomain; i++ {
			share.DPFKey = randomizeDPFKey(share.DPFKey)
			start := time.Now()
			kl.ExpandDPF(share)
			totalTime += time.Since(start).Microseconds()
		}

	} else if fssType == FSSRange {
		// (blackbox) inequality is 2*logn invocations of a DPF
		for i := uint(0); i < 2*fssDomain; i++ {
			share.DPFKey = randomizeDPFKey(share.DPFKey)
			start := time.Now()
			kl.ExpandDPF(share)
			totalTime += time.Since(start).Microseconds()
		}
	}

	start := time.Now()
	kl.Audit(share)
	totalTime += time.Since(start).Microseconds()

	return totalTime
}

func benchmarkPACLSymmetricKeyVFSS(
	kl *paclsposs.KeyList,
	share *paclsposs.ProofShare,
	fssDomain uint,
	fssType int) int64 {

	totalTime := int64(0)

	if fssType == FSSEquality {
		// equality is just DPF
		share.DPFKey = randomizeDPFKey(share.DPFKey)
		start := time.Now()
		kl.ExpandVDPF(share)
		totalTime += time.Since(start).Microseconds()

	} else if fssType == FSSInequality {

		// (blackbox) inequality is logn invocations of a DPF
		for i := uint(0); i < fssDomain; i++ {
			share.DPFKey = randomizeDPFKey(share.DPFKey)
			start := time.Now()
			kl.ExpandVDPF(share)
			totalTime += time.Since(start).Microseconds()
		}

	} else if fssType == FSSRange {
		// (blackbox) inequality is 2*logn invocations of a DPF
		for i := uint(0); i < 2*fssDomain; i++ {
			share.DPFKey = randomizeDPFKey(share.DPFKey)
			start := time.Now()
			kl.ExpandVDPF(share)
			totalTime += time.Since(start).Microseconds()
		}
	}

	klsk, _, _ := paclsk.GenerateBenchmarkKeyList(kl.NumKeys, kl.FSSDomain, paclsk.Inclusion, kl.NumKeys)

	start := time.Now()
	klsk.Audit(
		&paclsk.ProofShare{
			DPFKey:      share.DPFKey,
			PrfKey:      share.PrfKey,
			ShareNumber: 0,
			KeyShare:    paclsk.NewEmptySlot(klsk.StatSecurity / 8),
		})
	totalTime += time.Since(start).Microseconds()

	return totalTime
}

func benchmarkPACLPublicKeyVFSS(
	kl *paclsposs.KeyList,
	share *paclsposs.ProofShare,
	fssDomain uint,
	fssType int) int64 {

	totalTime := int64(0)

	if fssType == FSSEquality {
		// equality is just DPF
		share.DPFKey = randomizeDPFKey(share.DPFKey)

	} else if fssType == FSSInequality {

		// (blackbox) inequality is logn invocations of a DPF
		for i := uint(0); i < fssDomain; i++ {
			share.DPFKey = randomizeDPFKey(share.DPFKey)
			start := time.Now()
			kl.ExpandVDPF(share)
			totalTime += time.Since(start).Microseconds()
		}

	} else if fssType == FSSRange {
		// (blackbox) inequality is 2*logn invocations of a DPF
		for i := uint(0); i < 2*fssDomain; i++ {
			share.DPFKey = randomizeDPFKey(share.DPFKey)
			start := time.Now()
			kl.ExpandVDPF(share)
			totalTime += time.Since(start).Microseconds()
		}
	}

	start := time.Now()
	kl.Audit(share)
	totalTime += time.Since(start).Microseconds()

	return totalTime
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
