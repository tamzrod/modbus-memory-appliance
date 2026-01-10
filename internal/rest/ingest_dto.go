package rest

type ingestRequest struct {
	Memory  string   `json:"memory"`
	Area    string   `json:"area"`
	Address uint16   `json:"address"`
	Bools   []int    `json:"bools,omitempty"`
	Values  []uint16 `json:"values,omitempty"`
}
