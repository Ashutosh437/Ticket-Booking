package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"ticket-system/internal/models"
)

var ErrTicketNotFound = errors.New("ticket not found")

type TicketRepository interface {
	CreateTicket(ctx context.Context, userID, title, description string) (*models.Ticket, error)
	GetTicketsByUserID(ctx context.Context, userID string) ([]*models.Ticket, error)
	GetTicketByIDAndUserID(ctx context.Context, ticketID, userID string) (*models.Ticket, error)
	UpdateTicketStatus(ctx context.Context, ticketID, userID string, newStatus models.TicketStatus) (*models.Ticket, error)
}

type ticketRepo struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) TicketRepository {
	return &ticketRepo{db: db}
}

func (r *ticketRepo) CreateTicket(ctx context.Context, userID, title, description string) (*models.Ticket, error) {
	now := time.Now().UTC()
	ticket := &models.Ticket{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Status:      models.StatusOpen,
		UserID:      userID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	query := `INSERT INTO tickets (id, title, description, status, user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, ticket.ID, ticket.Title, ticket.Description, string(ticket.Status), ticket.UserID, ticket.CreatedAt, ticket.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("could not create ticket: %w", err)
	}

	return ticket, nil
}

func (r *ticketRepo) GetTicketsByUserID(ctx context.Context, userID string) ([]*models.Ticket, error) {
	query := `SELECT id, title, description, status, user_id, created_at, updated_at FROM tickets WHERE user_id = ? ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("could not query tickets: %w", err)
	}
	defer rows.Close()

	tickets := make([]*models.Ticket, 0)
	for rows.Next() {
		var t models.Ticket
		var statusStr string
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &statusStr, &t.UserID, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("could not scan ticket: %w", err)
		}
		t.Status = models.TicketStatus(statusStr)
		tickets = append(tickets, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row error: %w", err)
	}

	return tickets, nil
}

func (r *ticketRepo) GetTicketByIDAndUserID(ctx context.Context, ticketID, userID string) (*models.Ticket, error) {
	query := `SELECT id, title, description, status, user_id, created_at, updated_at FROM tickets WHERE id = ? AND user_id = ?`
	row := r.db.QueryRowContext(ctx, query, ticketID, userID)

	var t models.Ticket
	var statusStr string
	err := row.Scan(&t.ID, &t.Title, &t.Description, &statusStr, &t.UserID, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTicketNotFound
		}
		return nil, fmt.Errorf("could not fetch ticket: %w", err)
	}

	t.Status = models.TicketStatus(statusStr)
	return &t, nil
}

func (r *ticketRepo) UpdateTicketStatus(ctx context.Context, ticketID, userID string, newStatus models.TicketStatus) (*models.Ticket, error) {
	// 1. Fetch current ticket to verify existence and check ownership & status
	ticket, err := r.GetTicketByIDAndUserID(ctx, ticketID, userID)
	if err != nil {
		return nil, err
	}

	// 2. Validate status transition rule
	if err := models.ValidateStatusTransition(ticket.Status, newStatus); err != nil {
		return nil, err
	}

	// 3. Perform update if status changed
	now := time.Now().UTC()
	query := `UPDATE tickets SET status = ?, updated_at = ? WHERE id = ? AND user_id = ?`
	_, err = r.db.ExecContext(ctx, query, string(newStatus), now, ticketID, userID)
	if err != nil {
		return nil, fmt.Errorf("could not update ticket status: %w", err)
	}

	ticket.Status = newStatus
	ticket.UpdatedAt = now
	return ticket, nil
}
