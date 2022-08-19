package main

type Experiment struct {
	NumKeys    uint64  `json:"num_keys"`
	AuthTimeMS []int64 `json:"pacl_auth_time_ms"`
}
