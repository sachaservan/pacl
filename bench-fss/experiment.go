package main

type Experiment struct {
	NumKeys                         uint64  `json:"num_keys"`
	FSSDomain                       uint64  `json:"fss_domain"`
	EqualityBaselineProcessingNS    []int64 `json:"equality_baseline_processing_ns"`
	EqualityBaselineVerProcessingNS []int64 `json:"equality_baseline_ver_processing_ns"`
	EqualityDPFPACLProcessingNS     []int64 `json:"equality_dpf_pacl_processing_ns"`
	EqualityDPFSKPACLProcessingNS   []int64 `json:"equality_dpf_sk_pacl_processing_ns"`
	EqualityVDPFPACLProcessingNS    []int64 `json:"equality_vdpf_pacl_processing_ns"`
	EqualityVDPFSKPACLProcessingNS  []int64 `json:"equality_vdpf_sk_pacl_processing_ns"`

	InequalityBaselineProcessingNS    []int64 `json:"inequality_baseline_processing_ns"`
	InequalityBaselineVerProcessingNS []int64 `json:"inequality_baseline_ver_processing_ns"`
	InequalityDPFPACLProcessingNS     []int64 `json:"inequality_dpf_pacl_processing_ns"`
	InequalityDPFSKPACLProcessingNS   []int64 `json:"inequality_dpf_sk_pacl_processing_ns"`
	InequalityVDPFPACLProcessingNS    []int64 `json:"inequality_vdpf_pacl_processing_ns"`
	InequalityVDPFSKPACLProcessingNS  []int64 `json:"inequality_vdpf_sk_pacl_processing_ns"`

	RangeBaselineProcessingNS    []int64 `json:"range_baseline_processing_ns"`
	RangeBaselineVerProcessingNS []int64 `json:"range_baseline_ver_processing_ns"`
	RangeDPFPACLProcessingNS     []int64 `json:"range_dpf_pacl_processing_ns"`
	RangeDPFSKPACLProcessingNS   []int64 `json:"range_dpf_sk_pacl_processing_ns"`
	RangeVDPFPACLProcessingNS    []int64 `json:"range_vdpf_pacl_processing_ns"`
	RangeVDPFSKPACLProcessingNS  []int64 `json:"range_vdpf_sk_pacl_processing_ns"`
}
