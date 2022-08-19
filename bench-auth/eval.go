package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"time"

	paclsposs "github.com/sachaservan/pacl/pacl-sposs"
)

func main() {

	n := []int{262144, 524288, 1048576, 2097152, 8388608}

	for i, numAccounts := range n {

		experiment := &Experiment{
			NumKeys: uint64(numAccounts),
		}

		experiment.AuthTimeMS = make([]int64, 0)

		for trial := 0; trial < 10; trial++ {
			experiment.AuthTimeMS = append(experiment.AuthTimeMS, benchmarkPACLAuthTime(numAccounts))
			fmt.Printf("Finished trial %v of %v\n", trial, 10)
		}

		fmt.Printf("Auth time @ %v accounts: %v\n", numAccounts, experiment.AuthTimeMS)

		experimentJSON, _ := json.MarshalIndent(experiment, "", " ")
		ioutil.WriteFile("experiment"+fmt.Sprint(i)+".json", experimentJSON, 0644)
	}

}
func benchmarkPACLAuthTime(numAccount int) int64 {

	// setup parameters
	group := paclsposs.DefaultGroup()
	n := uint(math.Log2(float64(numAccount)))
	kl, key, idx := paclsposs.GenerateBenchmarkKeyList(
		uint64(numAccount), n, group, paclsposs.Equality, 0)

	// client-side computation (precomputed here because we're
	// benchmarking the server overhead).
	shares := kl.NewProof(idx, key)

	auditB := kl.Audit(shares[0])

	start := time.Now()
	auditA := kl.Audit(shares[0])
	kl.CheckAudit(auditA, auditB)

	return time.Since(start).Milliseconds()
}
