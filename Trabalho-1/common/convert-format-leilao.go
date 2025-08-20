package common

func LeilaoToString(leilao Leilao) []string {
	var leilaoString []string

	leilaoString = append(leilaoString, leilao.ID)
	leilaoString = append(leilaoString, leilao.description)
	leilaoString = append(leilaoString, leilao.startDate.String())
	leilaoString = append(leilaoString, leilao.endDate.String())

	return leilaoString
}

func LeiloesToCsv(leiloes []Leilao) [][]string {
	var data [][]string
	var leilaoString []string

	for i, leilao := range leiloes {
		leilaoString = LeilaoToString(leilao)
		data = append(data, leilaoString)
		print(i)
	}

	return data
}
