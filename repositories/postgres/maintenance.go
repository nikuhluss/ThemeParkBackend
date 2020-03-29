package postgres

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

var selectMaintenance = psql.Select("rides_maintenance.*").From("rides maintenance")

// MaintenanceRepository implements the MaintenanceRepository interface for postgres.
type MaintenanceRepository struct {
	db *sqlx.DB
}

// NewMaintenanceRepository creates a new MaintenanceRepository instance using the given
// database instance.
func NewMaintenanceRepository(db *sqlx.DB) *MaintenanceRepository {
	return &MaintenanceRepository{db}
}

// GetByID fetches a maintenance from the database using the given ID.
func (rr *MaintenanceRepository) GetByID(ID string) (*models.Maintenance, error) {
	db := rr.db
	udb := db.Unsafe()

	query, _ := selectMaintenance.Where(sq.Eq{"Maintenance.ID": ID}).MustSql()

	rides_maintenance := models.Maintenance{}
	err := udb.Get(&rides_maintenance, query, ID)
	if err != nil {
		return nil, err
	}

	return &rides_maintenance, nil
}

// Fetch fetches all maintenance from the database.
func (rr *MaintenanceRepository) Fetch() ([]*models.Maintenance, error) {
	db := rr.db
	udb := db.Unsafe()

	query, _ := selectMaintenance.MustSql()

	maintenance := []*models.Maintenance{}
	err := udb.Select(&maintenance, query)
	if err != nil {
		return nil, err
	}

	return maintenance, err
}

// Store creates an entry for the given maintenance model in the database.
func (rr *MaintenanceRepository) Store(maintenance *models.Maintenance) error {
	db := rr.db

	insertMaintenance, _, _ := psql.
		Insert("rides_maintenance").
		Columns("ID", "ride_id", "maintenance_type_id", "description", "cost", "start_datetime", "end_datetime").
		Values("?", "?", "?", "?", "?", "?", "?").
		ToSql()

	_, err := db.Exec(insertMaintenance, maintenance.ID, maintenance.RideID, maintenance.MaintenanceTypeID, maintenance.Description, maintenance.Cost, maintenance.StartDateTime, maintenance.EndDateTime)
	if err != nil {
		return fmt.Errorf("inserMaintenance: %s", err)
	}

	return nil
}

// Update updates an existing entry in the database for the given Maintenance model.
func (rr *MaintenanceRepository) Update(maintenance *models.Maintenance) error {
	db := rr.db

	updateMaintenance, _, _ := psql.
		Update("rides_maintenance").
		Set("ride_id", "?").
		Set("maintenance_type_id", "?").
		Set("description", "?").
		Set("start_datetime", "?").
		Set("end_datetime", "?").
		Where("id = ?").
		ToSql()

	_, err := db.Exec(updateMaintenance, maintenance.RideID, maintenance.MaintenanceTypeID, maintenance.Description, maintenance.Cost, maintenance.StartDateTime, maintenance.EndDateTime)
	if err != nil {
		return fmt.Errorf("updateMaintenance: %s", err)
	}

	return nil
}

// Delete deletes an existing entry in the database for the given Maintenance ID.
func (rr *MaintenanceRepository) Delete(ID string) error {
	db := rr.db

	deleteMaintenance, _, _ := psql.Delete("rides_maintenance").Where("id = ?").ToSql()

	_, err := db.Exec(deleteMaintenance, ID)
	if err != nil {
		return fmt.Errorf("deleteMaintenance: %s", err)
	}

	return nil
}