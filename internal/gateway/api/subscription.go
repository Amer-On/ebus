package api

import (
	"context"
	"ebus/internal/domain"
	"encoding/json"
	"fmt"
	"net/http"
)

func (a *API) subscribeHandler(w http.ResponseWriter, r *http.Request) {
	subscribtion := new(domain.Subscription)
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(subscribtion)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode JSON: %v", err), http.StatusBadRequest)
		return
	}
	err = a.subscriptionService.Subscribe(*subscribtion)

	if err != nil {
		http.Error(w, fmt.Sprintf("Service error: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) publishHandler(w http.ResponseWriter, r *http.Request) {
	a.logger.Info("Start publishing")
	event := new(domain.RawEvent)
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(event)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode JSON: %v", err), http.StatusBadRequest)
		return
	}

	err = a.publicationService.Publish(context.Background(), event.Topic, *event)
	if err != nil {
		http.Error(w, fmt.Sprintf("Service error: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
