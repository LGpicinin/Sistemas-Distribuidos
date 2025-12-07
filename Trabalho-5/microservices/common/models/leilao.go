package models

import (
	common "common/utils"
	"encoding/json"
	"time"
)

type Leilao struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

func CreateLeilao(ID string, description string, startDate time.Time, endDate time.Time) Leilao {
	var leilao Leilao = Leilao{
		ID:          ID,
		Description: description,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	return leilao
}

func (leilao *Leilao) ToByteArray() []byte {
	leilaoByteArray, err := json.Marshal(*leilao)
	if err != nil {
		common.FailOnError(err, "Erro ao converter leilao para []byte")
	}

	return leilaoByteArray
}

func (leilao *Leilao) FromByteArray(byteArray []byte) {
	err := json.Unmarshal(byteArray, leilao)
	common.FailOnError(err, "Erro ao converter []byte para leilao")
}

func (leilao *Leilao) HasStarted() bool {
	return time.Now().Compare(leilao.StartDate) >= 0
}

func (leilao *Leilao) HasEnded() bool {
	return time.Now().Compare(leilao.EndDate) >= 0
}

func (leilao *Leilao) Print() string {
	return "Leilão:\n" +
		"\tID: " + leilao.ID + "\n" +
		"\tDescrição: " + leilao.Description + "\n" +
		"\tData de início: " + leilao.StartDate.String() + "\n" +
		"\tData de finalização: " + leilao.EndDate.String() + "\n"
}

type ByStartDate []Leilao

func (a ByStartDate) Len() int           { return len(a) }
func (a ByStartDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByStartDate) Less(i, j int) bool { return a[i].StartDate.Compare(a[j].StartDate) == -1 }

type ByEndDate []Leilao

func (a ByEndDate) Len() int           { return len(a) }
func (a ByEndDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByEndDate) Less(i, j int) bool { return a[i].EndDate.Compare(a[j].EndDate) == -1 }
