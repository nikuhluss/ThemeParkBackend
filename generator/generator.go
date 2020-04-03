package generator

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

// Rand wraps the standard rand.Rand.
type Rand struct {
	*rand.Rand
}

// FromStringSlice picks a random element from the given string slice.
func (r *Rand) FromStringSlice(s []string) (string, error) {
	if len(s) <= 0 {
		return "", fmt.Errorf("slice must be non-empty")
	}
	return s[r.Intn(len(s))], nil
}

// Inserter struct is useful for inserting mock values in the database.
type Inserter struct {
	execer Execer
	rand   *Rand
}

// NewInserter creates a new Inserter instance with a
// deterministic rand.
func NewInserter(execer Execer) *Inserter {
	source := rand.NewSource(0)
	rand := &Rand{rand.New(source)}

	return &Inserter{
		execer,
		rand,
	}
}

// DoInsert starts inserting everything in the database.
func (i *Inserter) DoInsert() error {

	// Genders

	fmt.Println("Inserting genders...")
	i.execer.Exec("TRUNCATE TABLE genders CASCADE")

	genderMale, err := InsertGender(i.execer, "Male")
	if err != nil {
		return err
	}

	genderFemale, err := InsertGender(i.execer, "Female")
	if err != nil {
		return err
	}

	genderOther, err := InsertGender(i.execer, "Other")
	if err != nil {
		return err
	}

	genders := []string{genderMale, genderFemale, genderOther}

	// Roles

	fmt.Println("Inserting roles...")
	i.execer.Exec("TRUNCATE TABLE roles CASCADE")

	roleWorker, err := InsertRole(i.execer, "Worker")
	if err != nil {
		return err
	}

	roleSupervisor, err := InsertRole(i.execer, "Supervisor")
	if err != nil {
		return err
	}

	roles := []string{roleWorker, roleSupervisor}

	// Maintenance Types

	fmt.Println("Inserting maintenance types...")
	i.execer.Exec("TRUNCATE TABLE maintenance_types CASCADE")

	maintenanceTypeTuneUp, err := InsertMaintenanceType(i.execer, "Tune Up")
	if err != nil {
		return err
	}

	maintenanceTypeReplacement, err := InsertMaintenanceType(i.execer, "Replacement")
	if err != nil {
		return err
	}

	maintenanceTypeFixed, err := InsertMaintenanceType(i.execer, "Fixed")
	if err != nil {
		return err
	}

	maintenanceTypes := []string{maintenanceTypeTuneUp, maintenanceTypeReplacement, maintenanceTypeFixed}

	// Customers

	fmt.Println("Inserting customers...")
	i.execer.Exec("TRUNCATE TABLE users CASCADE")

	totalCustomers := 1000
	customers := make([]string, 0, totalCustomers)
	for idx := 0; idx < totalCustomers; idx++ {
		username := fmt.Sprintf("customer%d", idx)
		email := fmt.Sprintf("%s@email.com", username)
		genderID, err := i.rand.FromStringSlice(genders)
		if err != nil {
			return err
		}

		customerID, err := InsertCustomer(i.execer, username, email)
		if err != nil {
			return err
		}

		err = InsertUserDetails(i.execer, customerID, genderID)
		if err != nil {
			return err
		}

		customers = append(customers, customerID)
	}

	// Employees

	fmt.Println("Inserting employees...")

	totalEmployees := 100
	employees := make([]string, 0, totalEmployees)
	for idx := 0; idx < totalEmployees; idx++ {
		username := fmt.Sprintf("employee%d", idx)
		email := fmt.Sprintf("%s@email.com", username)
		roleID, err := i.rand.FromStringSlice(roles)
		if err != nil {
			return err
		}
		genderID, err := i.rand.FromStringSlice(genders)
		if err != nil {
			return err
		}

		employeeID, err := InsertEmployee(i.execer, username, email, roleID)
		if err != nil {
			return err
		}

		err = InsertUserDetails(i.execer, employeeID, genderID)
		if err != nil {
			return err
		}

		employees = append(employees, employeeID)
	}

	// Rides

	fmt.Println("Inserting rides...")
	i.execer.Exec("TRUNCATE TABLE rides CASCADE")

	totalRides := 10
	rides := make([]string, 0, totalRides)
	for idx := 0; idx < totalRides; idx++ {
		rideID, err := InsertRide(i.execer)
		if err != nil {
			return err
		}

		rides = append(rides, rideID)
	}

	// Maintenance

	fmt.Println("Inserting maintenance jobs...")
	i.execer.Exec("TRUNCATE TABLE rides_maintenance CASCADE")

	allMaintenance := make([]string, 0)

	for _, rideID := range rides {
		maintenance, err := i.doInsertMaintenance(rideID, maintenanceTypes)
		if err != nil {
			return err
		}

		allMaintenance = append(allMaintenance, maintenance...)
	}

	// Tickets

	// fmt.Println("Inserting tickets...")
	// i.execer.Exec("TRUNCATE TABLE rides_maintenance CASCADE")

	for _, _ = range rides {
	}

	return nil
}

func (i *Inserter) doInsertMaintenance(rideID string, maintenanceTypes []string) ([]string, error) {

	allMaintenance := make([]string, 0)

	start := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)
	monthDuration := time.Hour * 24 * 7 * 4
	yearDuration := monthDuration * 12

	// assuming:
	// 1 year = 12 months
	// 1 month = 4 weeks
	// 1 week = 7 days
	// 1 day = 24 hours
	for year := 0; year < 10; year++ {
		for month := 0; month < 12; month++ {

			monthStart := start.Add(yearDuration * time.Duration(year)).Add(monthDuration * time.Duration(month))
			maintenancePerMonth := i.rand.Intn(10)

			for idx := 0; idx < maintenancePerMonth; idx++ {

				maintenanceStart := monthStart.Add(monthDuration / time.Duration(maintenancePerMonth))
				maintenanceEnd := sql.NullTime{Time: maintenanceStart.Add(time.Hour * 24), Valid: true}

				maintenanceType, err := i.rand.FromStringSlice(maintenanceTypes)
				if err != nil {
					return nil, err
				}

				maintenanceID, err := InsertMaintenanceWithStartAndEnd(i.execer, rideID, maintenanceType, maintenanceStart, maintenanceEnd)
				if err != nil {
					return nil, err
				}

				allMaintenance = append(allMaintenance, maintenanceID)

			} // maintenance
		} // month
	} // year

	return allMaintenance, nil
}
