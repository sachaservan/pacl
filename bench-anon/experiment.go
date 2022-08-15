package main

type Experiment struct {
	NumKeys               uint64  `json:"num_accounts"`
	ServerExpressMS       []int64 `json:"server_express_ms"`
	ServerSpectrumMS      []int64 `json:"server_spectrum_ms"`
	ServerExpressCloakMS  []int64 `json:"server_express_cloak_ms"`
	ServerSpectrumCloakMS []int64 `json:"server_spectrum_cloak_ms"`
}
