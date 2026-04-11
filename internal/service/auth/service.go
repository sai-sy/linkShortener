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
var ErrInvalidCredentials = errors.New("invalid credentials")

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

func (s *Service) Register(ctx context.Context, w http.ResponseWriter, email, password, firstname string) (db.AuthUser, error) {
	trimmedEmail := strings.TrimSpace(strings.ToLower(email))
	trimmedFirstname := strings.TrimSpace(firstname)
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

	profileFirstname := pgtype.Text{Valid: false}
	if trimmedFirstname != "" {
		profileFirstname = pgtype.Text{String: trimmedFirstname, Valid: true}
	}

	if _, err := s.db.CreateProfile(ctx, db.CreateProfileParams{
		UserID:    user.ID,
		Firstname: profileFirstname,
		Surname:   pgtype.Text{Valid: false},
	}); err != nil {
		return db.AuthUser{}, err
	}

	workspaceName := trimmedFirstname
	if workspaceName == "" {
		workspaceName = trimmedEmail
	}
	workspaceName = workspaceName + "'s Workspace"

	workspace, err := s.db.CreateWorkspace(ctx, workspaceName)
	if err != nil {
		return db.AuthUser{}, err
	}

	if err := s.db.CreateWorkspaceMember(ctx, db.CreateWorkspaceMemberParams{
		WorkspaceID: workspace.ID,
		ProfileID:   user.ID,
		Role:        "owner",
	}); err != nil {
		return db.AuthUser{}, err
	}

	ownerPermissions := []string{
		"workspace:read",
		"workspace:update",
		"routemap:create",
		"routemap:read",
		"routemap:update",
		"routemap:delete",
	}
	memberPermissions := []string{
		"workspace:read",
		"routemap:read",
	}

	for _, permission := range ownerPermissions {
		if err := s.db.CreateRolePermission(ctx, db.CreateRolePermissionParams{
			WorkspaceID: workspace.ID,
			Role:        "owner",
			Permission:  permission,
		}); err != nil {
			return db.AuthUser{}, err
		}
	}

	for _, permission := range memberPermissions {
		if err := s.db.CreateRolePermission(ctx, db.CreateRolePermissionParams{
			WorkspaceID: workspace.ID,
			Role:        "member",
			Permission:  permission,
		}); err != nil {
			return db.AuthUser{}, err
		}
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

func (s *Service) Login(ctx context.Context, w http.ResponseWriter, email, password string) (db.AuthUser, error) {
	trimmedEmail := strings.TrimSpace(strings.ToLower(email))
	if trimmedEmail == "" || password == "" {
		return db.AuthUser{}, ErrInvalidCredentials
	}

	user, err := s.db.GetAuthUserByEmail(ctx, trimmedEmail)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.AuthUser{}, ErrInvalidCredentials
		}
		return db.AuthUser{}, err
	}

	if !s.checkPasswordHash(password, user.PasswordHash) {
		return db.AuthUser{}, ErrInvalidCredentials
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
