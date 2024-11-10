package handlers

import (
	"encoding/json"
	"main/models"
	"main/services"
	"net/http"

	"github.com/gorilla/mux"
)

type AssetHandler struct {
	service *services.AssetService
}

func NewAssetHandler(service *services.AssetService) *AssetHandler {
	return &AssetHandler{service: service}
}

func (h *AssetHandler) CreateAssetHandler(w http.ResponseWriter, r *http.Request) {
	var novoAsset models.Asset
	if err := json.NewDecoder(r.Body).Decode(&novoAsset); err != nil {
		responseJSON(w, http.StatusBadRequest, AppError{Message: ErrInvalidData, Code: http.StatusBadRequest})
		return
	}

	ctx := r.Context()
	userID, ok := ctx.Value(userIDKey).(string)
	if !ok {
		responseJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user ID"})
		return
	}

	createAsset, err := h.service.CreateAsset(ctx, userID, novoAsset)
	if err != nil {
		responseJSON(w, http.StatusInternalServerError, AppError{Message: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	responseJSON(w, http.StatusCreated, createAsset)
}

func (h *AssetHandler) UpdateAssetHandler(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]

	var atualizadoAsset models.Asset
	if err := json.NewDecoder(r.Body).Decode(&atualizadoAsset); err != nil {
		responseJSON(w, http.StatusBadRequest, AppError{Message: ErrInvalidData, Code: http.StatusBadRequest})
		return
	}

	ctx := r.Context()

	userID, ok := ctx.Value(userIDKey).(string)
	if !ok {
		responseJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user ID"})
		return
	}

	assetAtualizado, err := h.service.UpdateAsset(ctx, userID, id, atualizadoAsset)
	if err != nil {
		responseJSON(w, http.StatusNotFound, AppError{Message: err.Error(), Code: http.StatusNotFound})
		return
	}

	responseJSON(w, http.StatusOK, assetAtualizado)
}

func (h *AssetHandler) DeleteAssetHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	ctx := r.Context()

	userID, ok := ctx.Value(userIDKey).(string)
	if !ok {
		responseJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user ID"})
		return
	}

	err := h.service.DeleteAsset(ctx, userID, id)
	if err != nil {
		responseJSON(w, http.StatusNotFound, AppError{Message: err.Error(), Code: http.StatusNotFound})
		return
	}

	responseJSON(w, http.StatusNoContent, nil)
}

func (h *AssetHandler) GetAssetsHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	userID, ok := ctx.Value(userIDKey).(string)
	if !ok {
		responseJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user ID"})
		return
	}

	assets, err := h.service.GetAssets(ctx, userID)
	if err != nil {
		responseJSON(w, http.StatusInternalServerError, AppError{Message: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	responseJSON(w, http.StatusOK, assets)
}
