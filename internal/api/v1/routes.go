package apiv1

import "github.com/sai-sy/linkShortener/internal/api/v1/handlers"

type Handler = handlers.Handler

func NewHandler() *handlers.Handler {
	return handlers.NewHandler()
}
