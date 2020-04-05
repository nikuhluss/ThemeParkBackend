package models

import (
	"database/sql"
	"time"
)

// Ticket struct contains information about a ticket.
type Ticket struct {
	ID                string    `json:"id"`
	UserID            string    `db:"user_id" json:"userId"`
	IsKid             bool      `db:"is_kid" json:"isKid"`
	PurchasePrice     float64   `db:"purchase_price" json:"purchasePrice"`
	PurchasedOn       time.Time `db:"purchased_on" json:"purchasedOn"`
	PurchaseReference string    `db:"purchase_reference" json:"purchaseReference"`

	Email     string         `json:"email"`
	FirstName sql.NullString `db:"first_name" json:"firstName"`
	LastName  sql.NullString `db:"last_name" json:"lastName"`
}

// TicketScan struct contains information about a ticket scan.
type TicketScan struct {
	ID       string    `json:"id"`
	TicketID string    `db:"ticket_id" json:"ticketId"`
	RideID   string    `db:"ride_id" json:"rideId"`
	ScanOn   time.Time `db:"scan_datetime" json:"scanOn"`

	UserID    string         `db:"user_id" json:"userId"`
	Email     string         `json:"email"`
	FirstName sql.NullString `db:"first_name" json:"firstName"`
	LastName  sql.NullString `db:"last_name" json:"lastName"`
}
