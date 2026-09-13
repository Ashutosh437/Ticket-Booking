package service

import (
	"context"
	"errors"
	"strings"

	"ticket-system/internal/models"
	"ticket-system/internal/repository"
)

type TicketService interface {
	CreateTicket(ctx context.Context, userID string, req models.CreateTicketRequest) (*models.Ticket, error)
	GetTickets(ctx context.Context, userID string) ([]*models.Ticket, error)
	GetTicketByID(ctx context.Context, ticketID, userID string) (*models.Ticket, error)
	UpdateTicketStatus(ctx context.Context, ticketID, userID string, req models.UpdateTicketStatusRequest) (*models.Ticket, error)
}

type ticketService struct {
	ticketRepo repository.TicketRepository
}

func NewTicketService(ticketRepo repository.TicketRepository) TicketService {
	return &ticketService{
		ticketRepo: ticketRepo,
	}
}

func (s *ticketService) CreateTicket(ctx context.Context, userID string, req models.CreateTicketRequest) (*models.Ticket, error) {
	if strings.TrimSpace(req.Title) == "" {
		return nil, errors.New("title is required")
	}
	if strings.TrimSpace(req.Description) == "" {
		return nil, errors.New("description is required")
	}

	return s.ticketRepo.CreateTicket(ctx, userID, req.Title, req.Description)
}

func (s *ticketService) GetTickets(ctx context.Context, userID string) ([]*models.Ticket, error) {
	return s.ticketRepo.GetTicketsByUserID(ctx, userID)
}

func (s *ticketService) GetTicketByID(ctx context.Context, ticketID, userID string) (*models.Ticket, error) {
	if strings.TrimSpace(ticketID) == "" {
		return nil, errors.New("ticket ID is required")
	}
	return s.ticketRepo.GetTicketByIDAndUserID(ctx, ticketID, userID)
}

func (s *ticketService) UpdateTicketStatus(ctx context.Context, ticketID, userID string, req models.UpdateTicketStatusRequest) (*models.Ticket, error) {
	if strings.TrimSpace(ticketID) == "" {
		return nil, errors.New("ticket ID is required")
	}

	return s.ticketRepo.UpdateTicketStatus(ctx, ticketID, userID, req.Status)
}
