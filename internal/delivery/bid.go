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
)

type BidHandler struct {
	service services.Bid
}

func NewBidHandler(service services.Bid) *BidHandler {
	return &BidHandler{
		service: service,
	}
}

func (h *BidHandler) GetUserBids(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := r.URL.Query()
	filterParams := operation.BidParams{}
	err := filterParams.Scan(params.Get("limit"), params.Get("offset"))
	if err != nil {
		operation.BadRequest(w)
		return
	}
	creator := params.Get("username")
	if creator == "" {
		operation.Unauthorized(w)
		return
	}
	bids, err := h.service.GetUserBids(r.Context(), filterParams, creator)
	if err != nil {
		operation.InternalServerError(w)
		return
	}
	response, err := json.Marshal(bids)
	if err != nil {
		operation.InternalServerError(w)
		return
	}
	w.Write(response)
}

func (h *BidHandler) GetBidsForTender(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := r.URL.Query()
	filterParams := operation.BidParams{}
	vars := mux.Vars(r)
	id := vars["tenderId"]
	if id == "" {
		operation.BadRequest(w)
		return
	}
	err := filterParams.Scan(params.Get("limit"), params.Get("offset"))
	if err != nil {
		operation.BadRequest(w)
		return
	}
	creator := params.Get("username")
	if creator == "" {
		operation.Unauthorized(w)
		return
	}
	bids, err := h.service.GetBidsForTender(r.Context(), id, filterParams)
	if err != nil {
		operation.InternalServerError(w)
		return
	}
	response, err := json.Marshal(bids)
	if err != nil {
		operation.InternalServerError(w)
		return
	}
	w.Write(response)
}

func (h *BidHandler) CreateBid(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		operation.BadRequest(w)
		return
	}
	defer r.Body.Close()

	b := entities.Bid{}
	err = json.Unmarshal(body, &b)
	if err != nil {
		operation.BadRequest(w)
		return
	}
	bid, err := h.service.CreateBid(r.Context(), b)
	if err != nil {
		operation.InternalServerError(w)
		return
	}
	response, err := json.Marshal(bid)
	if err != nil {
		operation.InternalServerError(w)
		return
	}
	w.Write(response)
}

func (h *BidHandler) GetBidStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := r.URL.Query()
	creator := params.Get("username")
	if creator == "" {
		operation.Unauthorized(w)
		return
	}
	vars := mux.Vars(r)
	id := vars["bidId"]
	if id == "" {
		operation.BadRequest(w)
		return
	}
	status, err := h.service.GetBidStatus(r.Context(), id)
	if err != nil {
		operation.InternalServerError(w)
		return
	}
	response, err := json.Marshal(status)
	if err != nil {
		operation.InternalServerError(w)
		return
	}
	w.Write(response)
}

func (h *BidHandler) ChangeBidStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := r.URL.Query()
	creator := params.Get("username")
	if creator == "" {
		operation.Unauthorized(w)
		return
	}
	str := params.Get("status")
	var status entities.BidStatus
	status.Scan(str)
	if status == "" {
		operation.BadRequest(w)
		return
	}
	vars := mux.Vars(r)
	id := vars["bidId"]
	if id == "" {
		operation.BadRequest(w)
		return
	}
	bid, err := h.service.ChangeBidStatus(r.Context(), status, id)
	if err != nil {
		operation.InternalServerError(w)
		return
	}
	response, err := json.Marshal(bid)
	if err != nil {
		operation.InternalServerError(w)
		return
	}
	w.Write(response)
}

func (h *BidHandler) SubmitBid(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := r.URL.Query()
	creator := params.Get("username")
	if creator == "" {
		operation.Unauthorized(w)
		return
	}
	decision := params.Get("decision")
	var status entities.BidStatus
	status.Scan(decision)
	if status == "" {
		operation.BadRequest(w)
		return
	}
	vars := mux.Vars(r)
	id := vars["bidId"]
	if id == "" {
		operation.BadRequest(w)
		return
	}
	bid, err := h.service.SubmitBid(r.Context(), status, id, creator)
	if err != nil {
		if errors.Is(err, services.ErrNoAccess) || errors.Is(err, services.ErrNotResponsible) {
			operation.Forbidden(w)
			return
		}
		operation.InternalServerError(w)
		return
	}
	response, err := json.Marshal(bid)
	if err != nil {
		operation.InternalServerError(w)
		return
	}
	w.Write(response)
}

func (h *BidHandler) EditBid(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := r.URL.Query()
	creator := params.Get("username")
	if creator == "" {
		operation.Unauthorized(w)
		return
	}
	vars := mux.Vars(r)
	id := vars["bidId"]
	if id == "" {
		operation.BadRequest(w)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		operation.BadRequest(w)
		return
	}
	defer r.Body.Close()

	t := entities.Bid{}
	err = json.Unmarshal(body, &t)
	if err != nil {
		operation.BadRequest(w)
		return
	}
	bid, err := h.service.EditBid(r.Context(), t, id)
	if err != nil {
		operation.InternalServerError(w)
		return
	}
	response, err := json.Marshal(bid)
	if err != nil {
		operation.InternalServerError(w)
		return
	}
	w.Write(response)
}
