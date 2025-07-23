package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type SubscriptionHandler struct {
	subscriptionRepo SubscriptionRepository
}

func NewSubscriptionHandler(subscriptionRepo SubscriptionRepository) *SubscriptionHandler {
	return &SubscriptionHandler{subscriptionRepo: subscriptionRepo}
}

func (h *SubscriptionHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/subscription", h.createSubscription).Methods("POST")
	router.HandleFunc("/subscription", h.updateSubscription).Methods("PUT")
	router.HandleFunc("/subscription", h.listSubscriptions).Methods("GET")
	router.HandleFunc("/subscription/summary", h.getSummary).Methods("GET")
	router.HandleFunc("/subscription/{id}", h.getSubscription).Methods("GET")
	router.HandleFunc("/subscription/{id}", h.deleteSubscription).Methods("DELETE")
}

// @Summary Create Subscription
// @Description Create subscription
// @ID create
// @Accept json
// @Success 201
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Param createData body Subscription true "Create data"
// @Router /subscription [post]
func (h *SubscriptionHandler) createSubscription(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	fmt.Println("Request body:", string(body))
	r.Body = io.NopCloser(bytes.NewBuffer(body)) // Восстанавливаем тело для повторного чтения

	var request Subscription
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ResponseError(w, http.StatusBadRequest, "invalid request format")
		return
	}

	if err := h.subscriptionRepo.Create(&request); err != nil {
		ResponseError(w, http.StatusInternalServerError, "failed to create subscription")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// @Summary Update Subscription
// @Description Update subscription data
// @ID update
// @Accept json
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Param updateData body Subscription true "Update data"
// @Router /subscription [put]
func (h *SubscriptionHandler) updateSubscription(w http.ResponseWriter, r *http.Request) {
	var request Subscription
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ResponseError(w, http.StatusBadRequest, "invalid request format")
		return
	}

	if err := h.subscriptionRepo.Update(&request); err != nil {
		ResponseError(w, http.StatusInternalServerError, "failed to update subscription")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Get subscription
// @Description Get subscription by id
// @ID get
// @Produce  json
// @Success 200 {object} Subscription
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Param id path int true "Subscription ID"
// @Router /subscription/{id} [get]
func (h *SubscriptionHandler) getSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subscriptionID, err := strconv.Atoi(vars["id"])
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, "invalid subscription ID")
		return
	}

	subscription, err := h.subscriptionRepo.Read(uint(subscriptionID))
	if err != nil {
		ResponseError(w, http.StatusNotFound, "no such subscription")
		return
	}

	ResponseJSON(w, http.StatusOK, subscription)
}

// @Summary Delete subscription
// @Description Delete subscription by id
// @ID delete
// @Produce  json
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Param id path int true "Subscription ID"
// @Router /subscription/{id} [delete]
func (h *SubscriptionHandler) deleteSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subscriptionID, err := strconv.Atoi(vars["id"])
	if err != nil {
		ResponseError(w, http.StatusBadRequest, "invalid subscription ID")
		return
	}

	if err := h.subscriptionRepo.Delete(uint(subscriptionID)); err != nil {
		ResponseError(w, http.StatusInternalServerError, "failed to delete subscription")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary List subscriptions
// @Description List all subscriptions
// @ID list
// @Produce json
// @Success 200 {object} []Subscription
// @Failure 500 {object} ErrorResponse
// @Router /subscription [get]
func (h *SubscriptionHandler) listSubscriptions(w http.ResponseWriter, r *http.Request) {
	subs, err := h.subscriptionRepo.List()
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, "failed to list subscriptions")
		return
	}

	ResponseJSON(w, http.StatusOK, subs)
}

// @Summary Get subscriptions cost
// @Description Get total cost of subscriptions for given filters
// @ID cost
// @Produce json
// @Success 200 {object} SummaryResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Param start_date query string false "Start Date"
// @Param end_date query string false "End Date"
// @Param user_id query string false "User ID"
// @Param service_name query string false "Service Name"
// @Router /subscription/summary [get]
func (h *SubscriptionHandler) getSummary(w http.ResponseWriter, r *http.Request) {
	var err error // kind of a hack

	var startDate time.Time
	if startDateQuery := r.URL.Query().Get("start_date"); startDateQuery != "" {
		startDate, err = time.Parse("01-2006", r.URL.Query().Get("start_date"))
		if err != nil {
			ResponseError(w, http.StatusBadRequest, "failed to parse start date")
			return
		}
	}

	var endDate time.Time
	if endDateQuery := r.URL.Query().Get("end_date"); endDateQuery != "" {
		endDate, err = time.Parse("01-2006", endDateQuery)
		if err != nil {
			ResponseError(w, http.StatusBadRequest, "failed to parse end date")
			return
		}
	}

	var userID uuid.UUID
	if userIDQuery := r.URL.Query().Get("user_id"); userIDQuery != "" {
		userID, err = uuid.Parse(userIDQuery)
		if err != nil {
			ResponseError(w, http.StatusBadRequest, "failed to parse user ID")
			return
		}
	}

	serviceName := r.URL.Query().Get("service_name")

	subs, err := h.subscriptionRepo.Filter(startDate, endDate, userID, serviceName)
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, "failed to get subscriptions")
		return
	}

	// It may be a good idea to put this logic in the DB request itself,
	// but the current Filter() seems to be a more universal tool, not for calculating the total cost only
	var cost uint = 0
	for _, sub := range subs {
		cost += sub.Price
	}

	ResponseJSON(w, http.StatusOK, SummaryResponse{Cost: cost})
}

type SummaryResponse struct {
	Cost uint
}
