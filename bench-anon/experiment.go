package main

type Experiment struct {
	NumKeys              uint64  `json:"num_keys"`
	ServerExpressMS      []int64 `json:"server_express_ms"`
	ServerSpectrumMS     []int64 `json:"server_spectrum_ms"`
	ServerExpressPACLMS  []int64 `json:"server_express_pacl_ms"`
	ServerSpectrumPACLMS []int64 `json:"server_spectrum_pacl_ms"`
}
