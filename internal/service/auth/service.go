package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sai-sy/linkShortener/internal/db"
	"golang.org/x/crypto/bcrypt"
)

var ErrEmailTaken = errors.New("email already in use")
var ErrInvalidSession = errors.New("invalid session")
var ErrInvalidCSRF = errors.New("invalid csrf token")

type Service struct {
	db *db.Queries
}

func NewService(db *db.Queries) *Service {
	return &Service{db: db}
}

func (s *Service) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *Service) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return nil == err
}

func (s *Service) createSessionToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func (s *Service) Register(ctx context.Context, w http.ResponseWriter, email, password string) (db.AuthUser, error) {
	trimmedEmail := strings.TrimSpace(strings.ToLower(email))
	if trimmedEmail == "" || password == "" {
		return db.AuthUser{}, errors.New("email and password required")
	}

	if _, err := s.db.GetAuthUserByEmail(ctx, trimmedEmail); err == nil {
		return db.AuthUser{}, ErrEmailTaken
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return db.AuthUser{}, err
	}

	hashBytes, err := s.hashPassword(password)
	if err != nil {
		return db.AuthUser{}, err
	}

	userID := uuid.New()
	user, err := s.db.CreateAuthUser(ctx, db.CreateAuthUserParams{
		ID:           pgtype.UUID{Bytes: userID, Valid: true},
		Email:        trimmedEmail,
		PasswordHash: string(hashBytes),
	})
	if err != nil {
		return db.AuthUser{}, err
	}

	if _, err := s.db.CreateProfile(ctx, db.CreateProfileParams{
		UserID:    user.ID,
		Firstname: pgtype.Text{Valid: false},
		Surname:   pgtype.Text{Valid: false},
	}); err != nil {
		return db.AuthUser{}, err
	}

	sessionToken, err := s.createSessionToken()
	if err != nil {
		return db.AuthUser{}, err
	}

	csrfToken, err := s.createSessionToken()
	if err != nil {
		return db.AuthUser{}, err
	}

	if err := s.db.CreateAuthSession(ctx, db.CreateAuthSessionParams{
		UserID:       user.ID,
		SessionToken: sessionToken,
		CsrfToken:    csrfToken,
	}); err != nil {
		return db.AuthUser{}, err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	return user, nil
}

func (s *Service) Authenticate(ctx context.Context, r *http.Request) (db.AuthSession, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return db.AuthSession{}, ErrInvalidSession
	}

	session, err := s.db.GetAuthSessionByToken(ctx, cookie.Value)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.AuthSession{}, ErrInvalidSession
		}
		return db.AuthSession{}, err
	}

	csrfToken := r.Header.Get("X-CSRF-Token")
	if csrfToken == "" || csrfToken != session.CsrfToken {
		return db.AuthSession{}, ErrInvalidCSRF
	}

	return session, nil
}
