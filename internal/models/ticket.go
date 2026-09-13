package models

import (
	"errors"
	"time"
)

type TicketStatus string

const (
	StatusOpen       TicketStatus = "open"
	StatusInProgress TicketStatus = "in_progress"
	StatusClosed     TicketStatus = "closed"
)

func (s TicketStatus) IsValid() bool {
	switch s {
	case StatusOpen, StatusInProgress, StatusClosed:
		return true
	default:
		return false
	}
}

type Ticket struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      TicketStatus `json:"status"`
	UserID      string       `json:"user_id"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type CreateTicketRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTicketStatusRequest struct {
	Status TicketStatus `json:"status"`
}

var (
	ErrInvalidStatus           = errors.New("invalid status value: must be 'open', 'in_progress', or 'closed'")
	ErrClosedTicketReopen      = errors.New("a closed ticket cannot be reopened or changed")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
)

func ValidateStatusTransition(current, next TicketStatus) error {
	if !next.IsValid() {
		return ErrInvalidStatus
	}

	if current == StatusClosed {
		if next == StatusOpen || next == StatusInProgress {
			return ErrClosedTicketReopen
		}
	}

	if current == StatusInProgress && next == StatusOpen {
		return ErrInvalidStatusTransition
	}

	return nil
}
