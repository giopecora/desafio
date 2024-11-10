package handlers

import (
	"encoding/json"
	"main/models"
	"main/services"
	"net/http"

	"github.com/gorilla/mux"
)

type DebtHandler struct {
	service *services.DebtService
}

func NewDebtHandler(service *services.DebtService) *DebtHandler {
	return &DebtHandler{service: service}
}

func (h *DebtHandler) CreateDebtHandler(w http.ResponseWriter, r *http.Request) {

	var requestData struct {
		UserId string      `json:"user_id"`
		Debt   models.Debt `json:"debt"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		responseJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid data"})
		return
	}

	ctx := r.Context()

	debt, err := h.service.CreateDebt(ctx, requestData.UserId, requestData.Debt)
	if err != nil {
		responseJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	responseJSON(w, http.StatusCreated, debt)
}

func (h *DebtHandler) UpdateDebtHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var atualizadaDebt models.Debt
	if err := json.NewDecoder(r.Body).Decode(&atualizadaDebt); err != nil {
		responseJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid data"})
		return
	}

	ctx := r.Context()

	debt, err := h.service.UpdateDebt(ctx, id, atualizadaDebt)
	if err != nil {
		responseJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	responseJSON(w, http.StatusOK, debt)
}

func (h *DebtHandler) DeleteDebtHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	ctx := r.Context()

	err := h.service.DeleteDebt(ctx, id)
	if err != nil {
		responseJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	responseJSON(w, http.StatusNoContent, nil)
}

func (h *DebtHandler) GetDebtsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId := mux.Vars(r)["user_id"]

	userDebts, err := h.service.GetDebts(ctx, userId)
	if err != nil {
		responseJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	responseJSON(w, http.StatusOK, userDebts)
}
