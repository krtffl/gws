package webui

import "github.com/krtffl/get-well-soon/internal/domain"

type Service struct {
	rep domain.GWSRepo
}

func NewSvc(rep domain.GWSRepo) *Service {
	return &Service{
		rep: rep,
	}
}

func (svc *Service) List() ([]*domain.GWS, error) {
	return svc.rep.List()
}

func (svc *Service) Create(gws *domain.GWS) (*domain.GWS, error) {
	return svc.rep.Create(gws)
}
