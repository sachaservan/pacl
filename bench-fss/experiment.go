package main

type Experiment struct {
	NumKeys                       uint64  `json:"num_keys"`
	NumSubkeys                    uint64  `json:"num_subkeys"`
	FSSDomain                     uint64  `json:"fss_domain"`
	GroupExponentiation           uint64  `json:"group_exp_us"`
	EqualityBaselineProcessing    []int64 `json:"equality_baseline_processing_us"`
	EqualityBaselineVerProcessing []int64 `json:"equality_baseline_ver_processing_us"`
	EqualityDPFPACLProcessing     []int64 `json:"equality_dpf_pacl_processing_us"`
	EqualityDPFSKPACLProcessing   []int64 `json:"equality_dpf_sk_pacl_processing_us"`
	EqualityVDPFPACLProcessing    []int64 `json:"equality_vdpf_pacl_processing_us"`
	EqualityVDPFSKPACLProcessing  []int64 `json:"equality_vdpf_sk_pacl_processing_us"`

	RangeBaselineProcessing    []int64 `json:"range_baseline_processing_us"`
	RangeBaselineVerProcessing []int64 `json:"range_baseline_ver_processing_us"`
	RangeDPFPACLProcessing     []int64 `json:"range_dpf_pacl_processing_us"`
	RangeDPFSKPACLProcessing   []int64 `json:"range_dpf_sk_pacl_processing_us"`
	RangeVDPFPACLProcessing    []int64 `json:"range_vdpf_pacl_processing_us"`
	RangeVDPFSKPACLProcessing  []int64 `json:"range_vdpf_sk_pacl_processing_us"`
}
