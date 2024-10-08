package interfaces

import "time"

type Tabloid struct {
	Nome             string
	DtInicioVigencia time.Time
	DtFimVigencia    time.Time
	Ativo            bool
	DtCadastro       time.Time
	DtAlteracao      time.Time
	RegiaoID         int
}
