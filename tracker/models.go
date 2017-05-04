package tracker

import "time"

type Donation struct {
	ID            int64     `json:"id"`
	RefTypeID     int     `json:"-"`
	RefID         int64     `json:"-"`
	CharacterName string    `json:"characterName"`
	CharacterID   int64     `json:"characterId"`
	Amount        float64   `json:"amount"`
	Date          time.Time `json:"date"`
	Reason        string `json:"-"`
}

type Donations []Donation
