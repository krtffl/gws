package webui

import "github.com/krtffl/gws/internal/domain"

type Service struct {
	rep domain.GWSRepo
}

func NewSvc(rep domain.GWSRepo) *Service {
	return &Service{
		rep: rep,
	}
}

func (svc *Service) List(limit, offset uint) ([]*domain.GWS, error) {
	return svc.rep.List(limit, offset)
}

func (svc *Service) Create(gws *domain.GWS) (*domain.GWS, error) {
	return svc.rep.Create(gws)
}
