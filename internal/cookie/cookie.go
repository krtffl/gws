package cookie

import (
	"net/http"
	"time"
)

func (svc *Service) CreateCookie(session string) (*http.Cookie, error) {
	hashed, err := svc.cookie.Encode(svc.cfg.Name, session)
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     svc.cfg.Name,
		Value:    hashed,
		Expires:  time.Now().AddDate(0, 0, int(svc.cfg.Expiry)),
		Path:     "/",
		HttpOnly: true,
	}

	return cookie, nil
}

func (svc *Service) RemoveCookie() *http.Cookie {
	return &http.Cookie{
		Name:    svc.cfg.Name,
		Value:   "",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
}
