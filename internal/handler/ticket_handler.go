package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"ticket-system/internal/middleware"
	"ticket-system/internal/models"
	"ticket-system/internal/repository"
	"ticket-system/internal/service"
)

type TicketHandler struct {
	ticketService service.TicketService
}

func NewTicketHandler(ticketService service.TicketService) *TicketHandler {
	return &TicketHandler{
		ticketService: ticketService,
	}
}

func (h *TicketHandler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	var req models.CreateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	ticket, err := h.ticketService.CreateTicket(r.Context(), userID, req)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusCreated, ticket)
}

func (h *TicketHandler) GetTickets(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	tickets, err := h.ticketService.GetTickets(r.Context(), userID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, tickets)
}

func (h *TicketHandler) GetTicketByID(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	ticketID := chi.URLParam(r, "id")
	ticket, err := h.ticketService.GetTicketByID(r.Context(), ticketID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrTicketNotFound) {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": "Ticket not found"})
			return
		}
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, ticket)
}

func (h *TicketHandler) UpdateTicketStatus(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	ticketID := chi.URLParam(r, "id")
	var req models.UpdateTicketStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	ticket, err := h.ticketService.UpdateTicketStatus(r.Context(), ticketID, userID, req)
	if err != nil {
		if errors.Is(err, repository.ErrTicketNotFound) {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": "Ticket not found"})
			return
		}
		if errors.Is(err, models.ErrInvalidStatus) ||
			errors.Is(err, models.ErrClosedTicketReopen) ||
			errors.Is(err, models.ErrInvalidStatusTransition) {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, ticket)
}
