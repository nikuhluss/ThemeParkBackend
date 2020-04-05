package postgres_test

import (
	"fmt"
	"testing"
	"time"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/internal/testutil"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/generator"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func setupTestTickets(db *sqlx.DB) ([]string, []string, []string, []string) {

	tx := db.MustBegin()
	tx.MustExec("TRUNCATE TABLE users CASCADE")
	tx.MustExec("TRUNCATE TABLE rides CASCADE")

	// 3 customers
	customer0 := generator.MustInsertCustomer(tx, "customer0", "customer0@email.com")
	customer1 := generator.MustInsertCustomer(tx, "customer1", "customer1@email.com")
	customer2 := generator.MustInsertCustomer(tx, "customer2", "customer2@email.com")

	// tickets per customer = customer index
	ticket0 := generator.MustInsertTicket(tx, customer1, time.Now().UTC())
	ticket1 := generator.MustInsertTicket(tx, customer2, time.Now().UTC())
	ticket2 := generator.MustInsertTicket(tx, customer2, time.Now().UTC())

	// 3 rides
	ride0 := generator.MustInsertRide(tx)
	ride1 := generator.MustInsertRide(tx)
	ride2 := generator.MustInsertRide(tx)

	// scans per ride = ride index
	scan0 := generator.MustInsertTicketScan(tx, ticket0, ride1, time.Now().UTC())
	scan1 := generator.MustInsertTicketScan(tx, ticket1, ride2, time.Now().UTC())
	scan2 := generator.MustInsertTicketScan(tx, ticket2, ride2, time.Now().UTC())

	err := tx.Commit()
	if err != nil {
		panic(err)
	}

	return []string{customer0, customer1, customer2},
		[]string{ticket0, ticket1, ticket2},
		[]string{ride0, ride1, ride2},
		[]string{scan0, scan1, scan2}
}

func TestTicketGetByIDSucceeds(t *testing.T) {
	ticketRepository, db, teardown := testutil.MakeTicketRepositoryFixture()
	defer teardown()

	beforeSetupTime := time.Now().UTC()
	_, tickets, _, _ := setupTestTickets(db)

	for _, ticketID := range tickets {

		ticket, err := ticketRepository.GetByID(ticketID)
		if !assert.Nil(t, err) {
			t.FailNow()
		}

		fmt.Println(beforeSetupTime)
		fmt.Println(ticket.PurchasedOn)
		fmt.Println(beforeSetupTime.Before(ticket.PurchasedOn))

		assert.Equal(t, ticketID, ticket.ID)
		assert.NotEmpty(t, ticket.UserID)
		assert.False(t, ticket.IsKid)
		assert.Greater(t, ticket.PurchasePrice, float64(0))
		assert.True(t, beforeSetupTime.Before(ticket.PurchasedOn))
		assert.NotEmpty(t, ticket.PurchaseReference)
	}
}

func TestTicketFetchSucceeds(t *testing.T) {
	ticketRepository, db, teardown := testutil.MakeTicketRepositoryFixture()
	defer teardown()

	_, ticketIDs, _, _ := setupTestTickets(db)

	tickets, err := ticketRepository.Fetch()
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Len(t, tickets, len(ticketIDs))
}

func TestTicketFetchForUserSucceeds(t *testing.T) {
	ticketRepository, db, teardown := testutil.MakeTicketRepositoryFixture()
	defer teardown()

	userIDs, _, _, _ := setupTestTickets(db)

	for idx, userID := range userIDs {
		tickets, err := ticketRepository.FetchForUser(userID)
		if !assert.Nil(t, err) {
			t.FailNow()
		}

		assert.Len(t, tickets, idx)
	}
}

func TestTicketFetchScansSucceeds(t *testing.T) {
	ticketRepository, db, teardown := testutil.MakeTicketRepositoryFixture()
	defer teardown()

	_, _, _, scanIDs := setupTestTickets(db)

	scans, err := ticketRepository.Fetch()
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Len(t, scans, len(scanIDs))
}

func TestTicketFetchScansForRideSucceeds(t *testing.T) {
	ticketRepository, db, teardown := testutil.MakeTicketRepositoryFixture()
	defer teardown()

	_, _, rideIDs, _ := setupTestTickets(db)

	for idx, rideID := range rideIDs {
		scans, err := ticketRepository.FetchScansForRide(rideID)
		if !assert.Nil(t, err) {
			t.FailNow()
		}

		assert.Len(t, scans, idx)
	}
}

func TestTicketFetchScansForUserSucceeds(t *testing.T) {
	ticketRepository, db, teardown := testutil.MakeTicketRepositoryFixture()
	defer teardown()

	userIDs, _, _, _ := setupTestTickets(db)

	for idx, userID := range userIDs {
		scans, err := ticketRepository.FetchScansForUser(userID)
		if !assert.Nil(t, err) {
			t.FailNow()
		}

		assert.Len(t, scans, idx)
	}
}

func TestTicketStoreSucceeds(t *testing.T) {
	ticketRepository, db, teardown := testutil.MakeTicketRepositoryFixture()
	defer teardown()

	userIDs, _, _, _ := setupTestTickets(db)
	userID := userIDs[0]
	email := "customer0@email.com"

	expectedTicket := &models.Ticket{
		ID:                "some-ticket-id",
		UserID:            userID,
		IsKid:             true,
		PurchasePrice:     10,
		PurchasedOn:       time.Now().UTC(),
		PurchaseReference: "some-purchase-reference-id",
		Email:             email,
	}

	err := ticketRepository.Store(expectedTicket)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	ticket, err := ticketRepository.GetByID(expectedTicket.ID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Equal(t, expectedTicket, ticket)
}

func TestTicketUpdateSucceeds(t *testing.T) {
	ticketRepository, db, teardown := testutil.MakeTicketRepositoryFixture()
	defer teardown()

	_, ticketIDs, _, _ := setupTestTickets(db)
	ticketID := ticketIDs[0]

	expectedTicket, err := ticketRepository.GetByID(ticketID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	expectedTicket.IsKid = true
	expectedTicket.PurchasePrice = 10
	expectedTicket.PurchasedOn = time.Now().UTC()
	expectedTicket.PurchaseReference = "some-purchase-reference-id"

	err = ticketRepository.Update(expectedTicket)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	ticket, err := ticketRepository.GetByID(expectedTicket.ID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Equal(t, expectedTicket, ticket)
}

func TestTicketStoreScanSucceeds(t *testing.T) {
	ticketRepository, db, teardown := testutil.MakeTicketRepositoryFixture()
	defer teardown()

	_, ticketIDs, rideIDs, _ := setupTestTickets(db)
	ticketID := ticketIDs[0]
	rideID := rideIDs[0] // NOTE: ride0 has no scans

	expectedScan := &models.TicketScan{
		ID:       "some-ticket-scan-id",
		TicketID: ticketID,
		RideID:   rideID,
		ScanOn:   time.Now().UTC(),
	}

	err := ticketRepository.StoreScan(expectedScan)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	scans, err := ticketRepository.FetchScansForRide(rideID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Len(t, scans, 1)
}
