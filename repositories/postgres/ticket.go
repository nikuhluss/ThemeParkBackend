package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

var selectTickets = psql.
	Select(
		"tickets.*",
		"(DATE_TRUNC('day', tickets.purchased_on) = DATE_TRUNC('day', NOW())) AS is_valid",
		"users.email",
		"user_details.first_name",
		"user_details.last_name",
	).
	From("tickets").
	Join("users ON users.id = tickets.user_id").
	LeftJoin("user_details ON user_details.user_id = tickets.user_id").
	OrderBy("tickets.purchased_on DESC")

var selectTicketScans = psql.
	Select("scans.*", "users.email", "user_details.first_name", "user_details.last_name").
	From("tickets_on_rides AS scans").
	Join("tickets ON tickets.id = scans.ticket_id").
	Join("users ON users.id = tickets.user_id").
	LeftJoin("user_details ON user_details.user_id = tickets.user_id").
	OrderBy("scans.scan_datetime DESC")

// TicketRepository implements the TicketRepository interface for postgres.
type TicketRepository struct {
	db *sqlx.DB
}

// NewTicketRepository returns a new TicketRepository instance.
func NewTicketRepository(db *sqlx.DB) *TicketRepository {
	return &TicketRepository{db}
}

// GetByID fetches a ticket using the given ID.
func (tr *TicketRepository) GetByID(ID string) (*models.Ticket, error) {
	db := tr.db
	udb := db.Unsafe()

	query, _ := selectTickets.Where(sq.Eq{"tickets.id": "$1"}).MustSql()

	ticket := models.Ticket{}
	err := udb.Get(&ticket, query, ID)
	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

// Fetch fetches all tickets.
func (tr *TicketRepository) Fetch() ([]*models.Ticket, error) {
	db := tr.db
	udb := db.Unsafe()

	query, _ := selectTickets.MustSql()

	tickets := []*models.Ticket{}
	err := udb.Select(&tickets, query)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

// FetchForUser fetches all tickets for the given user.
func (tr *TicketRepository) FetchForUser(userID string) ([]*models.Ticket, error) {
	db := tr.db
	udb := db.Unsafe()

	query, _ := selectTickets.Where(sq.Eq{"tickets.user_id": "$1"}).MustSql()

	tickets := []*models.Ticket{}
	err := udb.Select(&tickets, query, userID)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

// FetchScans fetches all ticket scans.
func (tr *TicketRepository) FetchScans() ([]*models.TicketScan, error) {
	db := tr.db
	udb := db.Unsafe()

	query, _ := selectTicketScans.MustSql()

	scans := []*models.TicketScan{}
	err := udb.Select(&scans, query)
	if err != nil {
		return nil, err
	}

	return scans, nil
}

// FetchScansForRide fetches all ticket scans for the given ride.
func (tr *TicketRepository) FetchScansForRide(rideID string) ([]*models.TicketScan, error) {
	db := tr.db
	udb := db.Unsafe()

	query, _ := selectTicketScans.Where(sq.Eq{"scans.ride_id": "$1"}).MustSql()

	scans := []*models.TicketScan{}
	err := udb.Select(&scans, query, rideID)
	if err != nil {
		return nil, err
	}

	return scans, nil
}

// FetchScansForUser fetches all ticket scans for the given user.
func (tr *TicketRepository) FetchScansForUser(userID string) ([]*models.TicketScan, error) {
	db := tr.db
	udb := db.Unsafe()

	query, _ := selectTicketScans.Where(sq.Eq{"tickets.user_id": "$1"}).MustSql()

	scans := []*models.TicketScan{}
	err := udb.Select(&scans, query, userID)
	if err != nil {
		return nil, err
	}

	return scans, nil
}

// Store creates a new ticket.
func (tr *TicketRepository) Store(ticket *models.Ticket) error {
	db := tr.db

	query, _, _ := psql.
		Insert("tickets").
		Columns("id", "user_id", "is_kid", "purchase_price", "purchased_on", "purchase_reference").
		Values("$1", "$2", "$3", "$4", "$5", "$6").
		ToSql()

	_, err := db.Exec(query, ticket.ID, ticket.UserID, ticket.IsKid, ticket.PurchasePrice, ticket.PurchasedOn, ticket.PurchaseReference)
	if err != nil {
		return err
	}

	return nil
}

// Update updates an existing ticket.
func (tr *TicketRepository) Update(ticket *models.Ticket) error {
	db := tr.db

	query, _, _ := psql.
		Update("tickets").
		Set("user_id", "$1").
		Set("is_kid", "$2").
		Set("purchase_price", "$3").
		Set("purchased_on", "$4").
		Set("purchase_reference", "$5").
		Where(sq.Eq{"id": "$6"}).
		ToSql()

	_, err := db.Exec(query, ticket.UserID, ticket.IsKid, ticket.PurchasePrice, ticket.PurchasedOn, ticket.PurchaseReference, ticket.ID)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes an existing ticket.
func (tr *TicketRepository) Delete(ticketID string) error {
	db := tr.db

	query, _, _ := psql.
		Delete("tickets").
		Where(sq.Eq{"id": "$1"}).
		ToSql()

	_, err := db.Exec(query, ticketID)
	if err != nil {
		return err
	}

	return nil
}

// StoreScan creates a new ticket scan.
func (tr *TicketRepository) StoreScan(ticketScan *models.TicketScan) error {
	db := tr.db

	query, _, _ := psql.
		Insert("tickets_on_rides").
		Columns("id", "ride_id", "ticket_id", "scan_datetime").
		Values("$1", "$2", "$3", "$4").
		ToSql()

	_, err := db.Exec(query, ticketScan.ID, ticketScan.RideID, ticketScan.TicketID, ticketScan.ScanOn)
	if err != nil {
		return err
	}

	return nil
}

// UpdateScan updates an existing ticket scan.
func (tr *TicketRepository) UpdateScan(ticketScan *models.TicketScan) error {
	db := tr.db

	query, _, _ := psql.
		Update("tickets_on_rides").
		Set("ride_id", "$1").
		Set("ticket_id", "$2").
		Set("scan_datetime", "$3").
		Where(sq.Eq{"id": "$4"}).
		ToSql()

	_, err := db.Exec(query, ticketScan.RideID, ticketScan.TicketID, ticketScan.ScanOn, ticketScan.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteScan deletes an existing ticket scan.
func (tr *TicketRepository) DeleteScan(ticketScanID string) error {
	db := tr.db

	query, _, _ := psql.
		Delete("tickets_on_rides").
		Where(sq.Eq{"id": "$1"}).
		ToSql()

	_, err := db.Exec(query, ticketScanID)
	if err != nil {
		return err
	}

	return nil
}
