package common

import (
	"encoding/json"
)

type JustLink struct {
	Link string `json:"link"`
}

type Link struct {
	Link   string `json:"link"`
	UserID string `json:"clientId"`
}

func CreateLink(li string, userId string) Link {
	var link Link = Link{
		Link:   li,
		UserID: userId,
	}

	return link
}

func (link *Link) ToByteArray() []byte {
	linkByteArray, err := json.Marshal(*link)
	if err != nil {
		FailOnError(err, "Erro ao converter link para []byte")
	}

	return linkByteArray
}

func (link *Link) FromByteArray(byteArray []byte) {
	err := json.Unmarshal(byteArray, link)
	FailOnError(err, "Erro ao converter []byte para link")
}

func (link *Link) Print() string {
	return "Link:\n" +
		"\tID do usuário: " + link.UserID + "\n" +
		"\tLink: " + link.Link + "\n"
}
