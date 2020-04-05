package impl

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/internal/mathutil"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
	repos "gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/repositories"
)

var (
	errTicketExists        = fmt.Errorf("ticket with the given ID already exists")
	errTicketDoesNotExists = fmt.Errorf("ticket with the given ID does not exists")
)

// TicketUsecaseImpl implements the TicketUsecase interface.
type TicketUsecaseImpl struct {
	ticketRepo repos.TicketRepository
	rideRepo   repos.RideRepository
}

// NewTicketUsecaseImpl returns a new TicketUsecaseImpl instance.
func NewTicketUsecaseImpl(ticketRepo repos.TicketRepository, rideRepo repos.RideRepository) *TicketUsecaseImpl {
	return &TicketUsecaseImpl{ticketRepo, rideRepo}
}

// GetByID fetches a ticket with the given ID from the repository.
func (tu *TicketUsecaseImpl) GetByID(ctx context.Context, ID string) (*models.Ticket, error) {
	ticket, err := tu.ticketRepo.GetByID(ID)
	if err != nil {
		return nil, errTicketDoesNotExists
	}
	return ticket, nil
}

// Fetch fetches all the tickets.
func (tu *TicketUsecaseImpl) Fetch(ctx context.Context) ([]*models.Ticket, error) {
	return tu.ticketRepo.Fetch()
}

// FetchForUser fetches all the tickets for the given user.
func (tu *TicketUsecaseImpl) FetchForUser(ctx context.Context, userID string) ([]*models.Ticket, error) {
	return tu.ticketRepo.FetchForUser(userID)
}

// FetchScans fetches all scans.
func (tu *TicketUsecaseImpl) FetchScans(ctx context.Context) ([]*models.TicketScan, error) {
	return tu.ticketRepo.FetchScans()
}

// FetchScansForUser fetches all scans for the given user.
func (tu *TicketUsecaseImpl) FetchScansForUser(ctx context.Context, userID string) ([]*models.TicketScan, error) {
	return tu.ticketRepo.FetchScansForUser(userID)
}

// FetchScansForRide fetches all scans for the given ride.
func (tu *TicketUsecaseImpl) FetchScansForRide(ctx context.Context, rideID string) ([]*models.TicketScan, error) {
	return tu.ticketRepo.FetchScansForRide(rideID)
}

// Store creates a new Ticket.
func (tu *TicketUsecaseImpl) Store(ctx context.Context, ticket *models.Ticket) error {
	_, err := tu.ticketRepo.GetByID(ticket.ID)
	if err == nil {
		return errTicketExists
	}

	uuid, err := GenerateUUID()
	if err != nil {
		return err
	}

	ticket.ID = uuid
	cleanTicket(ticket)
	err = validateTicket(ticket)
	if err != nil {
		return err
	}

	err = tu.ticketRepo.Store(ticket)
	if err != nil {
		return err
	}

	return nil
}

// Update updates an existing ticket.
func (tu *TicketUsecaseImpl) Update(ctx context.Context, ticket *models.Ticket) error {
	_, err := tu.ticketRepo.GetByID(ticket.ID)
	if err != nil {
		return errTicketDoesNotExists
	}

	cleanTicket(ticket)
	err = validateTicket(ticket)
	if err != nil {
		return err
	}

	err = tu.ticketRepo.Update(ticket)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes an existing ticket.
func (tu *TicketUsecaseImpl) Delete(ctx context.Context, ID string) error {
	_, err := tu.ticketRepo.GetByID(ID)
	if err != nil {
		return errTicketDoesNotExists
	}

	err = tu.ticketRepo.Delete(ID)
	if err != nil {
		return err
	}

	return nil
}

// ScanTicket creates a new ticket scan and returns the created object.
func (tu *TicketUsecaseImpl) ScanTicket(ctx context.Context, ticketID string, rideID string) (*models.TicketScan, error) {

	_, err := tu.ticketRepo.GetByID(ticketID)
	if err != nil {
		return nil, errTicketDoesNotExists
	}

	_, err = tu.rideRepo.GetByID(rideID)
	if err != nil {
		return nil, errRideDoesNotExists
	}

	uuid, err := GenerateUUID()
	if err != nil {
		return nil, err
	}

	scan := models.TicketScan{
		ID:       uuid,
		TicketID: ticketID,
		RideID:   rideID,
		ScanOn:   time.Now().UTC(),
	}

	err = tu.ticketRepo.StoreScan(&scan)
	if err != nil {
		return nil, err
	}

	return &scan, nil
}

func cleanTicket(ticket *models.Ticket) {
	ticket.ID = strings.TrimSpace(ticket.ID)
	ticket.UserID = strings.TrimSpace(ticket.UserID)
	ticket.PurchasePrice = mathutil.ClampFloat64(ticket.PurchasePrice, 0, 1000)
	ticket.PurchaseReference = strings.TrimSpace(ticket.PurchaseReference)
}

func validateTicket(ticket *models.Ticket) error {
	if len(ticket.ID) <= 0 {
		return fmt.Errorf("validateTicket: ID must be non-empty")
	}

	if len(ticket.UserID) <= 0 {
		return fmt.Errorf("validateTicket: UserID must be non-empty")
	}

	if len(ticket.PurchaseReference) <= 0 {
		return fmt.Errorf("validateTicket: PurchaseReference must be non-empty")
	}

	return nil
}
