package api

import (
	"ebus/internal/gateway/service"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

type API struct {
	logger *zap.Logger
	router *chi.Mux

	subscriptionService *service.SubscribtionService
	publicationService  *service.PublicationService
}

func NewAPI(
	logger *zap.Logger,
	subscriptionService *service.SubscribtionService,
	publicationService *service.PublicationService,
) *API {
	api := &API{
		logger:              logger,
		subscriptionService: subscriptionService,
		publicationService:  publicationService,
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Post("/subscribe", api.subscribeHandler)
	router.Post("/publish", api.publishHandler)

	api.router = router
	return api
}

func (a *API) GetHandler() *chi.Mux { return a.router }
