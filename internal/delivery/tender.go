package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"zadanie-6105/internal/services"
)

type TenderHandler struct {
	service services.TenderService
}

func NewTenderHandler(service services.TenderService) TenderHandler {
	return TenderHandler{
		service: service,
	}
}

// offset and limit add
func (h *TenderHandler) GetTenderList(w http.ResponseWriter, r *http.Request) {
	tenders, err := h.service.GetTenderList(r.Context())
	if err != nil {
		fmt.Print(tenders)
		w.WriteHeader(500)
	}
	response, err := json.Marshal(tenders)
	if err != nil {
		w.WriteHeader(500)
	}
	w.Write(response)
}
