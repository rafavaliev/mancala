package lobby

import "gorm.io/gorm"

type Lobby struct {
	gorm.Model
	Slug       string `json:"slug" gorm:"uniqueIndex"`
	NumClients int    `json:"num_clients"`
}
