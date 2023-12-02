package cookie

import (
	"net/http"

	"github.com/gorilla/securecookie"

	"github.com/krtffl/gws/internal/config"
)

type Service struct {
	cfg    *config.Cookie
	cookie *securecookie.SecureCookie
}

func New(
	cfg *config.Cookie,
) *Service {
	cookie := securecookie.New(
		[]byte(cfg.HashKey), []byte(cfg.BlockKey))
	return &Service{
		cfg:    cfg,
		cookie: cookie,
	}
}

func (svc *Service) Retrieve(r *http.Request) (string, error) {
	c, err := r.Cookie(svc.cfg.Name)
	if err != nil {
		return "", ErrInvalidCookie(err.Error())
	}

	var unhashed string
	err = svc.cookie.Decode(svc.cfg.Name, c.Value, &unhashed)
	if err != nil {
		return "", ErrInvalidCookie(err.Error())
	}

	return unhashed, nil
}
