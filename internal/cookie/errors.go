package cookie

import (
	"fmt"

	"github.com/krtffl/gws/internal/domain"
)

func ErrInvalidCookie(details string) error {
	return fmt.Errorf("%s:  %s", domain.InvalidCookieError, details)
}
