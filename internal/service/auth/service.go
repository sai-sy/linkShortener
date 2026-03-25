package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sai-sy/linkShortener/internal/db"
	"golang.org/x/crypto/bcrypt"
)

var ErrEmailTaken = errors.New("email already in use")

type Service struct {
	db *db.Queries
}

func NewService(db *db.Queries) *Service {
	return &Service{db: db}
}

func (s *Service) Register(ctx context.Context, email, password string) (db.AuthUser, error) {
	trimmedEmail := strings.TrimSpace(strings.ToLower(email))
	if trimmedEmail == "" || password == "" {
		return db.AuthUser{}, errors.New("email and password required")
	}

	if _, err := s.db.GetAuthUserByEmail(ctx, trimmedEmail); err == nil {
		return db.AuthUser{}, ErrEmailTaken
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return db.AuthUser{}, err
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return db.AuthUser{}, err
	}

	userID := uuid.New()
	return s.db.CreateAuthUser(ctx, db.CreateAuthUserParams{
		ID:           pgtype.UUID{Bytes: userID, Valid: true},
		Email:        trimmedEmail,
		PasswordHash: string(hashBytes),
	})
}
