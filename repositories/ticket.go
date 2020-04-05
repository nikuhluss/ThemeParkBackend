package repositories

import (
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

// TicketRepository defines the interface for interacting with tickets and
// ticket scans.
type TicketRepository interface {
	GetByID(ID string) (*models.Ticket, error)

	Fetch() ([]*models.Ticket, error)
	FetchForUser(userID string) ([]*models.Ticket, error)

	FetchScans() ([]*models.TicketScan, error)
	FetchScansForRide(rideID string) ([]*models.TicketScan, error)
	FetchScansForUser(rideID string) ([]*models.TicketScan, error)

	Store(ticket *models.Ticket) error
	Update(ticket *models.Ticket) error
	Delete(ticketID string) error

	StoreScan(ticketScan *models.TicketScan) error
	UpdateScan(ticketScan *models.TicketScan) error
	DeleteScan(ticketScanID string) error
}
