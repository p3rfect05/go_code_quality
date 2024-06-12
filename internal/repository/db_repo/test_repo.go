package db_repo

import (
	"errors"
	"time"

	"github.com/p3rfect05/go_proj/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {

	var newID int

	return newID, nil
}

// InsertRoomRestriction inserts the room's restriction for particular period (from "start" to "end")
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	return nil
}

// SearchAvailabilityByDateByRoomID returns "true" if the room "roomID" is available for given date range,
// false otherwise
func (m *testDBRepo) SearchAvailabilityByDateByRoomID(start, end time.Time, roomID int) (bool, error) {

	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms if any for given date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil
}

// SearchRoomByID searches a room by ID
func (m *testDBRepo) SearchRoomByID(roomID int) (models.Room, error) {
	var room models.Room
	if roomID > 3 || roomID < 2 {
		return room, errors.New("some error")
	}

	return room, nil
}
func (m *testDBRepo) GetUserByID(ID int) (models.User, error) {
	var user models.User
	return user, nil
}
func (m *testDBRepo) AuthenticateUser(email, testPassword string) (int, string, error) {
	return 1, "", nil
}
func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}
