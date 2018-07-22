package xlog

// Result query result
type Result struct {
	Records []Record `json:"records"`
	Limit   int      `json:"limit"`
}
