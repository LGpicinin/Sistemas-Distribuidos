package common

import "time"

type Leilao struct {
	ID          string
	description string
	startDate   time.Time
	endDate     time.Time
}

func CreateLeilao(ID string, description string, startDate time.Time, endDate time.Time) Leilao {
	var leilao Leilao = Leilao{
		ID:          ID,
		description: description,
		startDate:   startDate,
		endDate:     endDate,
	}

	return leilao
}
