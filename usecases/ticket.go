package usecases

import (
	"context"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

// TicketUsecase is the usecase for interacting with tickets.
type TicketUsecase interface {
	GetByID(ctx context.Context, ID string) (*models.Ticket, error)

	Fetch(ctx context.Context) ([]*models.Ticket, error)
	FetchForUser(ctx context.Context, userID string) ([]*models.Ticket, error)

	FetchScans(ctx context.Context) ([]*models.TicketScan, error)
	FetchScansForUser(ctx context.Context, userID string) ([]*models.TicketScan, error)
	FetchScansForRide(ctx context.Context, rideID string) ([]*models.TicketScan, error)

	Store(ctx context.Context, ticket *models.Ticket) error
	Update(ctx context.Context, ticket *models.Ticket) error
	Delete(ctx context.Context, ID string) error

	ScanTicket(ctx context.Context, ticketID string, rideID string) (*models.TicketScan, error)
}
