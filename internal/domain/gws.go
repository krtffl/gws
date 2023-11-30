package domain

type GWS struct {
	Id      string `db:"Id"      json:"id"`
	From    string `db:"From"    json:"from"`
	Message string `db:"Message" json:"message"`
	Memory  []byte `db:"Memory"  json:"memory"`
}

type GWSRepo interface {
	List() ([]*GWS, error)
	Create(*GWS) (*GWS, error)
}
