package main

type Experiment struct {
	DBSize                    uint64  `json:"db_size"`
	ItemSize                  uint64  `json:"item_size"`
	ServerPIRProcessingMS     []int64 `json:"server_pir_processing_ms"`
	ServerPIRPACLProcessingMS []int64 `json:"server_pir_pacl_processing_ms"`
}
