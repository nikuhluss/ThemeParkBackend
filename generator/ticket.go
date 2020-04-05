package generator

import (
	"fmt"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v4"
)

// InsertTicket inserts a new ticket for the given user.
func InsertTicket(execer Execer, userID string, purchasedOn time.Time) (string, error) {
	ID := gofakeit.UUID()
	isKid := false
	purchasePrice := 50.0
	purchaseReference := gofakeit.UUID()

	insertTicketQuery := `
	INSERT INTO tickets (id, user_id, is_kid, purchase_price, purchased_on, purchase_reference)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := execer.Exec(insertTicketQuery, ID, userID, isKid, purchasePrice, purchasedOn, purchaseReference)
	if err != nil {
		return "", err
	}

	return ID, nil
}

// InsertTicketScan inserts a ticket scan using the given ticket ID, ride ID, and scan time.
func InsertTicketScan(execer Execer, ticketID, rideID string, scanOn time.Time) (string, error) {
	ID := gofakeit.UUID()

	insertScanQuery := `
	INSERT INTO tickets_on_rides (id, ride_id, ticket_id, scan_datetime)
	VALUES ($1, $2, $3, $4)
	`

	_, err := execer.Exec(insertScanQuery, ID, rideID, ticketID, scanOn)
	if err != nil {
		return "", err
	}

	return ID, nil
}

// BulkInsertTicket inserts ticket scans in bulk by combining the given user ID
// with the given purchase times.
// see: https://stackoverflow.com/a/25192138
func BulkInsertTicket(execer Execer, userID string, purchaseTimes []time.Time) ([]string, error) {

	// constructs the multi-valued insert query

	totalTickets := len(purchaseTimes)
	tickets := make([]string, 0, totalTickets)
	valueStrings := make([]string, 0, totalTickets)
	valueArgs := make([]interface{}, 0, totalTickets*4)

	for idx, purchasedOn := range purchaseTimes {
		ID := gofakeit.UUID()
		isKid := false
		purchasePrice := 50.0
		purchaseReference := gofakeit.UUID()

		tickets = append(tickets, ID)
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d)", idx*6+1, idx*6+2, idx*6+3, idx*6+4, idx*6+5, idx*6+6))
		valueArgs = append(valueArgs, ID, userID, isKid, purchasePrice, purchasedOn, purchaseReference)
	}

	insertMultipleScansQuery := fmt.Sprintf(`
	INSERT INTO tickets (id, user_id, is_kid, purchase_price, purchased_on, purchase_reference)
	VALUES %s
	`, strings.Join(valueStrings, ","))

	// executes query

	_, err := execer.Exec(insertMultipleScansQuery, valueArgs...)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

// BulkInsertTicketScan inserts ticket scans in bulk for the given ride. Note that ticketIDs
// and scanTimes are matched by index.
// see: https://stackoverflow.com/a/25192138
func BulkInsertTicketScan(execer Execer, rideID string, ticketIDs []string, scanTimes []time.Time) ([]string, error) {

	// totalScans is the minimum length between ticketIDs and scanTimes
	// (no missing info)

	totalScans := len(ticketIDs)
	if len(scanTimes) < totalScans {
		totalScans = len(scanTimes)
	}

	// constructs the multi-valued insert query

	ticketScans := make([]string, 0, totalScans)
	valueStrings := make([]string, 0, totalScans)
	valueArgs := make([]interface{}, 0, totalScans*4)

	for idx := 0; idx < totalScans; idx++ {
		ID := gofakeit.UUID()
		ticketScans = append(ticketScans, ID)
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d)", idx*4+1, idx*4+2, idx*4+3, idx*4+4))
		valueArgs = append(valueArgs, ID, rideID, ticketIDs[idx], scanTimes[idx])
	}

	insertMultipleScansQuery := fmt.Sprintf(`
	INSERT INTO tickets_on_rides (id, ride_id, ticket_id, scan_datetime)
	VALUES %s
	`, strings.Join(valueStrings, ","))

	// executes query

	_, err := execer.Exec(insertMultipleScansQuery, valueArgs...)
	if err != nil {
		return nil, err
	}

	return ticketScans, nil
}

// MustInsertTicket is like InsertTicket but panics on error.
func MustInsertTicket(mustExecer MustExecer, userID string, purchasedOn time.Time) string {
	return MustInsert(InsertTicket(&AsExecer{mustExecer}, userID, purchasedOn))
}

// MustInsertTicketScan is like InsertTicketScan but panics on error.
func MustInsertTicketScan(mustExecer MustExecer, ticketID, rideID string, scanOn time.Time) string {
	return MustInsert(InsertTicketScan(&AsExecer{mustExecer}, ticketID, rideID, scanOn))
}
