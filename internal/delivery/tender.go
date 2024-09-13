package delivery

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"zadanie-6105/internal/delivery/operation"
	"zadanie-6105/internal/repositories/entities"
	"zadanie-6105/internal/services"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type TenderHandler struct {
	service services.Tender
	logger  *zap.Logger
}

func NewTenderHandler(service services.Tender, logger *zap.Logger) *TenderHandler {
	return &TenderHandler{
		service: service,
		logger:  logger,
	}
}

func (h *TenderHandler) GetTenderList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := r.URL.Query()
	filterParams := operation.TenderListParams{}
	err := filterParams.Scan(params.Get("limit"), params.Get("offset"), params.Get("service_type"))
	if err != nil {
		h.logger.Error(err.Error())
		operation.BadRequest(w)
		return
	}
	tenders, err := h.service.GetTenderList(r.Context(), filterParams)
	if err != nil {
		h.logger.Error(err.Error())
		operation.InternalServerError(w)
		return
	}
	response, err := json.Marshal(tenders)
	if err != nil {
		h.logger.Error(err.Error())
		operation.WriteResponse(w, 500, []byte(err.Error()))
		return
	}
	w.Write(response)
}

func (h *TenderHandler) CreateTender(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err.Error())
		operation.BadRequest(w)
		return
	}
	defer r.Body.Close()

	t := entities.Tender{}
	err = json.Unmarshal(body, &t)
	if err != nil {
		h.logger.Error(err.Error())
		operation.BadRequest(w)
		return
	}
	tender, err := h.service.CreateTender(r.Context(), t)
	if err != nil {
		h.logger.Error(err.Error())
		operation.InternalServerError(w)
		return
	}
	response, err := json.Marshal(tender)
	if err != nil {
		h.logger.Error(err.Error())
		operation.InternalServerError(w)
		return
	}
	w.Write(response)
}

func (h *TenderHandler) GetTenderByUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := r.URL.Query()
	filterParams := operation.TenderListParams{}
	err := filterParams.Scan(params.Get("limit"), params.Get("offset"), "")
	if err != nil {
		h.logger.Error(err.Error())
		operation.BadRequest(w)
		return
	}
	creator := params.Get("username")
	if creator == "" {
		operation.Unauthorized(w)
		return
	}
	tenders, err := h.service.GetTenderByUser(r.Context(), creator, filterParams)
	if err != nil {
		h.logger.Error(err.Error())
		operation.InternalServerError(w)
		return
	}
	response, err := json.Marshal(tenders)
	if err != nil {
		h.logger.Error(err.Error())
		operation.InternalServerError(w)
		return
	}
	w.Write(response)
}

func (h *TenderHandler) GetTenderStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := r.URL.Query()
	creator := params.Get("username")
	if creator == "" {
		h.logger.Error("no username")
		operation.Unauthorized(w)
		return
	}
	vars := mux.Vars(r)
	id := vars["tenderId"]
	if id == "" {
		h.logger.Error("no id")
		operation.BadRequest(w)
		return
	}
	status, err := h.service.GetTenderStatus(r.Context(), id, creator)
	if err != nil {
		h.logger.Error(err.Error())
		operation.InternalServerError(w)
		return
	}
	response, err := json.Marshal(status)
	if err != nil {
		h.logger.Error(err.Error())
		operation.InternalServerError(w)
		return
	}
	w.Write(response)
}

func (h *TenderHandler) ChangeTenderStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := r.URL.Query()
	creator := params.Get("username")
	str := params.Get("status")
	var status entities.TenderStatus
	status.Scan(str)
	if status == "" {
		operation.BadRequest(w)
		return
	}
	if creator == "" || status == "" {
		operation.Unauthorized(w)
		return
	}
	vars := mux.Vars(r)
	id := vars["tenderId"]
	if id == "" {
		operation.BadRequest(w)
		return
	}
	tender, err := h.service.ChangeTenderStatus(r.Context(), status, id, creator)
	if err != nil {
		h.logger.Error(err.Error())
		if errors.Is(services.ErrNotResponsible, err) || errors.Is(services.ErrNoAccess, err) {
			operation.WriteResponse(w, 403, []byte(`{"reason": "`+err.Error()+`"}`))
			return
		}
		operation.InternalServerError(w)
		return
	}
	response, err := json.Marshal(tender)
	if err != nil {
		h.logger.Error(err.Error())
		operation.InternalServerError(w)
		return
	}
	w.Write(response)
}

func (h *TenderHandler) EditTender(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := r.URL.Query()
	creator := params.Get("username")
	if creator == "" {
		operation.Unauthorized(w)
		return
	}
	vars := mux.Vars(r)
	id := vars["tenderId"]
	if id == "" {
		operation.BadRequest(w)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err.Error())
		operation.BadRequest(w)
		return
	}
	defer r.Body.Close()

	t := entities.Tender{}
	err = json.Unmarshal(body, &t)
	if err != nil {
		h.logger.Error(err.Error())
		operation.BadRequest(w)
		return
	}
	tender, err := h.service.EditTender(r.Context(), t, id, creator)
	if err != nil {
		h.logger.Error(err.Error())
		operation.InternalServerError(w)
		return
	}
	response, err := json.Marshal(tender)
	if err != nil {
		h.logger.Error(err.Error())
		operation.InternalServerError(w)
		return
	}
	w.Write(response)
}
