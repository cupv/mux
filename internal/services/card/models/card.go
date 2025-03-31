package models

type Card struct {
    ID         string `json:"id"`
    Word       string `json:"word"`
    Meaning    string `json:"meaning"`
    Example    string `json:"example,omitempty"`
    CreatedAt  int64  `json:"created_at"`
}