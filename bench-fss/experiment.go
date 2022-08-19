package main

type Experiment struct {
	NumKeys                         uint64  `json:"num_keys"`
	NumSubkeys                      uint64  `json:"num_subkeys"`
	FSSDomain                       uint64  `json:"fss_domain"`
	GroupExponentiationNS           uint64  `json:"group_exp_ns"`
	EqualityBaselineProcessingNS    []int64 `json:"equality_baseline_processing_ns"`
	EqualityBaselineVerProcessingNS []int64 `json:"equality_baseline_ver_processing_ns"`
	EqualityDPFPACLProcessingNS     []int64 `json:"equality_dpf_pacl_processing_ns"`
	EqualityDPFSKPACLProcessingNS   []int64 `json:"equality_dpf_sk_pacl_processing_ns"`
	EqualityVDPFPACLProcessingNS    []int64 `json:"equality_vdpf_pacl_processing_ns"`
	EqualityVDPFSKPACLProcessingNS  []int64 `json:"equality_vdpf_sk_pacl_processing_ns"`

	RangeBaselineProcessingNS    []int64 `json:"range_baseline_processing_ns"`
	RangeBaselineVerProcessingNS []int64 `json:"range_baseline_ver_processing_ns"`
	RangeDPFPACLProcessingNS     []int64 `json:"range_dpf_pacl_processing_ns"`
	RangeDPFSKPACLProcessingNS   []int64 `json:"range_dpf_sk_pacl_processing_ns"`
	RangeVDPFPACLProcessingNS    []int64 `json:"range_vdpf_pacl_processing_ns"`
	RangeVDPFSKPACLProcessingNS  []int64 `json:"range_vdpf_sk_pacl_processing_ns"`
}
