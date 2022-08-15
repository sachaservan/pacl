package main

import (
	"crypto/elliptic"
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"time"

	"github.com/sachaservan/pacl/algebra"
	paclsk "github.com/sachaservan/pacl/pacl-sk"
	paclsposs "github.com/sachaservan/pacl/pacl-sposs"
	dpf "github.com/sachaservan/vdpf"
)

func main() {

	n := []int{512, 1024, 2048, 4096, 8192, 16384, 32768, 65536, 131072, 262144, 524288, 1048576, 2097152, 4194304}

	for i, numAccounts := range n {

		experiment := &Experiment{
			NumKeys: uint64(numAccounts),
		}
		experiment.ServerExpressMS = make([]int64, 0)
		experiment.ServerExpressCloakMS = make([]int64, 0)
		experiment.ServerSpectrumMS = make([]int64, 0)
		experiment.ServerSpectrumCloakMS = make([]int64, 0)

		for trial := 0; trial < 10; trial++ {
			experiment.ServerExpressMS = append(experiment.ServerExpressMS, benchmarkVanillaExpress(numAccounts))
			experiment.ServerExpressCloakMS = append(experiment.ServerExpressCloakMS, benchmarkPACLExpress(numAccounts))
			experiment.ServerSpectrumMS = append(experiment.ServerSpectrumMS, benchmarkVanillaSpectrum(numAccounts))
			experiment.ServerSpectrumCloakMS = append(experiment.ServerSpectrumCloakMS, benchmarkPACLSpectrum(numAccounts))
			fmt.Printf("Finished trial %v of %v\n", trial, 10)
		}

		fmt.Printf("Express          @ %v mailboxes: %v\n", numAccounts, experiment.ServerExpressMS[0])
		fmt.Printf("Express (PACL)   @ %v mailboxes: %v\n", numAccounts, experiment.ServerExpressCloakMS[0])
		fmt.Printf("Spectrum         @ %v mailboxes: %v\n", numAccounts, experiment.ServerSpectrumMS[0])
		fmt.Printf("Spectrum (PACL)  @ %v mailboxes: %v\n", numAccounts, experiment.ServerSpectrumCloakMS[0])

		experimentJSON, _ := json.MarshalIndent(experiment, "", " ")
		ioutil.WriteFile("experiment"+fmt.Sprint(i)+".json", experimentJSON, 0644)
	}

}

func benchmarkVanillaExpress(numMailboxes int) int64 {
	prfKey := dpf.GeneratePRFKey()
	client := dpf.ClientDPFInitialize(prfKey)

	// Express requires setting DPF domain to 128 bits for security
	keyA, _ := client.GenDPFKeys(12345, 128)
	server := dpf.ServerDPFInitialize(client.PrfKey)

	// Precompute the randomness used in Express to audit the DPF
	// evaluation. Note: expresses uses a short seed s to obtain
	// a random vector of field elements.
	// We assume this expansion takes negligible time and so we skip it
	// (it only makes Express *slower*) when compared to Express+PACL

	// Express uses p =  2^128 − 159
	p := big.NewInt(2)
	p.Exp(p, big.NewInt(128), nil)
	p.Sub(p, big.NewInt(159))
	expressField := algebra.NewField(p)

	// precompute the random values and DPF inputs.
	// Note: DPF evaluation time isn't dependent on input
	// just evaluate on x = 1...numMailboxes
	r := make([]*algebra.FieldElement, numMailboxes)
	x := make([]uint64, numMailboxes)
	for i := 0; i < numMailboxes; i++ {
		r[i] = expressField.RandomElement()
		x[i] = uint64(i)
	}

	start := time.Now()

	shares := server.BatchEval(keyA, x)

	// Express audit protocol takes the inner product between
	// the expanded DPF shares (in the field F_p) and the randomness
	// vectors r and R = r^2. In this loop we compute R
	R := make([]*algebra.FieldElement, numMailboxes)
	for i := 0; i < numMailboxes; i++ {
		R[i] = expressField.Mul(r[i], r[i])
	}
	for i := 0; i < numMailboxes; i++ {
		// we don't use the result; just for timing purposes
		el := expressField.NewElement(big.NewInt(int64(shares[i])))
		expressField.Mul(el, r[i])
		expressField.Mul(el, R[i])
	}

	return time.Since(start).Milliseconds()
}

func benchmarkPACLExpress(numMailboxes int) int64 {

	// setup parameters
	n := uint(math.Log2(float64(numMailboxes)))
	kl, key, _ := paclsk.GenerateBenchmarkKeyList(uint64(numMailboxes), n)

	shares := kl.NewProof(0, key)
	auditB := kl.Audit(shares[1])

	start := time.Now()

	// audit (includes VDPF expansion)
	auditA := kl.Audit(shares[0])
	kl.CheckAudit(auditA, auditB)

	return time.Since(start).Milliseconds()
}

func benchmarkVanillaSpectrum(numChannels int) int64 {

	curve := elliptic.P256()

	prfKey := dpf.GeneratePRFKey()
	client := dpf.ClientDPFInitialize(prfKey)

	// Spectrum requires the DPF output to be a PRG *seed* but it doesn't
	// need to have the DPF *domain* be 128. We therefore set it to
	// log(# channels) to ensure efficiency.
	bits := uint(math.Log2(float64(numChannels)))
	keyA, _ := client.GenDPFKeys(0, bits)
	server := dpf.ServerDPFInitialize(client.PrfKey)

	// To make sure a client has knowledge of the "channel key",
	// Spectrum takes an inner product between the DPF expansion shares and the channel keys
	// ``in the exponent'' of an EC exponent field (Spectrum requires the DPF
	// output to be in the EC exponent field).
	// Specifically each server has a list of ``public keys'' g^a1 .. g^an
	// where n is the number of channels and following the inner product obtains
	// a share g^[ai] assuming the DPF expanded to 1 in only a single location

	// Here we pre-compute the DPF output for simplicity since our DPF interface
	// is not designed to cast outputs into a specific Field (it can trivially
	// be made to support this functionality if required, however.)
	scalars := make([]*big.Int, numChannels)
	for i := 0; i < numChannels; i++ {
		_, scalars[i], _ = RandomCurveScalar(curve, crand.Reader)
	}

	start := time.Now()

	// unlike Express, Spectrum can benefit from a "full domain" evaluation
	// optimization since the servers expand the DPF on *all* indices 1...n
	// we ignore the output since we will be using the pre-computed scalars
	// from above.
	server.FullDomainEval(keyA)

	// once the DPF is evaluated we must take the
	// inner product of the shares with each key ai in the exponent
	// of the group. In the EC group this corresponds to a scalar mult
	// followed by addition which results in the value g^[ai]
	// where [.] denotes a secret-share
	A := &Point{curve.Params().Gx, curve.Params().Gy}
	for i := 0; i < numChannels; i++ {
		P := curvePointScalarMult(curve, scalars[i])
		A = curvePointAdd(curve, A, P)
	}

	return time.Since(start).Milliseconds()
}

func benchmarkPACLSpectrum(numChannels int) int64 {

	// setup parameters
	group := paclsposs.DefaultGroup()
	n := uint(math.Log2(float64(numChannels)))
	kl, key, idx := paclsposs.GenerateBenchmarkKeyList(uint64(numChannels), n, group)

	// client-side computation (precomputed here because we're
	// benchmarking the server overhead).
	shares := kl.NewProof(idx, key)

	auditB := kl.Audit(shares[0])

	start := time.Now()
	auditA := kl.Audit(shares[0])
	kl.CheckAudit(auditA, auditB)

	return time.Since(start).Milliseconds()
}
