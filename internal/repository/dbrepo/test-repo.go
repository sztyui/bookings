package dbrepo

import (
	"errors"
	"time"

	"github.com/sztyui/bookings/internal/models"
)

// AllUsers lists all of the users from the database
func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation to the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// if the room id is 2 then fail otherwise pass
	if res.RoomID == 2 {
		return 0, errors.New("some error")
	}
	return 1, nil
}

// InsertRoomRestriction is for inserting room restriction into the database.
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == 3 {
		return errors.New("some error")
	}
	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exists, and false if no availability exists
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	if roomID == 2 {
		return false, nil
	}
	if roomID == 3 {
		return false, errors.New("some error")
	}
	return true, nil
}

// SearchAvailabilityForAllRooms is for checking availability for the rooms if any for given date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

// GetRoomByID gets a Room by ID
func (m *testDBRepo) GetRoomByID(ID int) (models.Room, error) {
	room := models.Room{}
	if ID > 2 {
		return room, errors.New("Some error")
	}
	return room, nil
}

// GetUserByID gets user by ID
func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	return models.User{}, nil
}

// UpdateUser updates a user in the database
func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

// Authenticate authenticates the user
func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 0, "", nil
}
