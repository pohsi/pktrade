package entity

type Card struct {
	ID    string `json:"id"`
	Owner int    `json:"owner"`
	Type  int    `json:"type"`
}
