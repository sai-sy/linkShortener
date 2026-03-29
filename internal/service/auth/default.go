package auth

import "github.com/sai-sy/linkShortener/internal/db"

var DefaultService *Service

func SetDefaultService(db *db.Queries) {
	DefaultService = NewService(db)
}
