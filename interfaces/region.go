package interfaces

import "time"

type Region struct {
	ID           uint      `json:"id"`
	Nome         string    `json:"nome"`
	Dt_cadastro  time.Time `json:"dt_cadastro"`
	Dt_alteracao time.Time `json:"dt_alteracao"`
}
