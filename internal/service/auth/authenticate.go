package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/sai-sy/linkShortener/internal/db"
)

func (s *Service) Authenticate(ctx context.Context, r *http.Request) (db.GetAuthSessionByTokenRow, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return db.GetAuthSessionByTokenRow{}, ErrInvalidSession
	}

	session, err := s.db.GetAuthSessionByToken(ctx, cookie.Value)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.GetAuthSessionByTokenRow{}, ErrInvalidSession
		}
		return db.GetAuthSessionByTokenRow{}, err
	}

	if r.Method != http.MethodGet && r.Method != http.MethodHead && r.Method != http.MethodOptions {
		csrfToken := r.Header.Get("X-CSRF-Token")
		if csrfToken == "" {
			if err := r.ParseForm(); err == nil {
				csrfToken = r.FormValue("csrf_token")
			}
		}
		if csrfToken == "" || csrfToken != session.CsrfToken {
			return db.GetAuthSessionByTokenRow{}, ErrInvalidCSRF
		}
	}

	return session, nil
}
