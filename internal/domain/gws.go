package domain

import "time"

type GWS struct {
	Id        string    `db:"Id"        json:"id"`
	From      string    `db:"From"      json:"from"`
	Message   string    `db:"Message"   json:"message"`
	Memory    []byte    `db:"Memory"    json:"memory"`
	CreatedAt time.Time `db:"CreatedAt" json:"createdAt"`
}

type GWSRepo interface {
	List(uint, uint) ([]*GWS, error)
	Create(*GWS) (*GWS, error)
}
