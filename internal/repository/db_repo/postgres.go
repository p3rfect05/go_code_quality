package db_repo

import (
	"context"
	"errors"
	"time"

	"github.com/p3rfect05/go_proj/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	var newID int
	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, 
                          				room_id, created_at, updated_at)
                          				values($1, $2, $3, $4, $5, $6, $7, $8, $9)
                          				returning id`
	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// InsertRoomRestriction inserts the room's restriction for particular period (from "start" to "end")
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id,
                               				created_at, updated_at, restriction_id)
                               				values($1, $2, $3, $4, $5, $6, $7)`
	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		time.Now(),
		time.Now(),
		r.RestrictionID)
	if err != nil {
		return err
	}
	return nil
}

// SearchAvailabilityByDateByRoomID returns "true" if the room "roomID" is available for given date range,
// false otherwise
func (m *postgresDBRepo) SearchAvailabilityByDateByRoomID(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	SELECT COUNT(id)
	FROM room_restrictions
	
	WHERE room_id = $1 AND $2 < end_date AND $3 > start_date
			`

	var numRows int

	err := m.DB.QueryRowContext(ctx, stmt, roomID, start, end).Scan(&numRows)
	if err != nil {
		return false, err
	}
	return numRows == 0, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms if any for given date range
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
	SELECT r.id, r.room_name
	FROM rooms AS r 
	WHERE r.id NOT IN (SELECT room_id 
	                   FROM room_restrictions AS rr 
	                   WHERE $1 < end_date AND $2 > start_date)
	`
	var rooms []models.Room
	rows, err := m.DB.QueryContext(ctx, stmt, start, end)
	if err != nil {
		return rooms, nil
	}
	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return nil, err // don't return any rooms if there is a mistake
		}
		rooms = append(rooms, room)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return rooms, nil
}

// SearchRoomByID searches a room by ID
func (m *postgresDBRepo) SearchRoomByID(roomID int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var room models.Room

	stmt := `
	SELECT id, room_name, created_at, updated_at 
	FROM rooms
	WHERE id = $1
		`
	row := m.DB.QueryRowContext(ctx, stmt, roomID)
	err := row.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return room, err
	}
	return room, nil
}

// GetUserByID returns a user by ID
func (m *postgresDBRepo) GetUserByID(ID int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User
	stmt := `
	SELECT id, first_name, last_name, email, password, access_level,
				created_at, updated_at
	FROM users WHERE id = $1
			`
	row := m.DB.QueryRowContext(ctx, stmt, ID)
	err := row.Scan(&user.FirstName, &user.LastName, &user.Email, &user.Password,
		&user.AccessLevel, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}
	return user, nil
}

// UpdateUser a user in a database
func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
	UPDATE users
	SET first_name = $1, last_name = $2, email = $3, access_level = $4, updated_at = $5
		`
	_, err := m.DB.ExecContext(ctx, stmt,
		u.FirstName, u.LastName, u.Email, u.AccessLevel, time.Now())

	if err != nil {
		return err
	}
	return nil
}

// AuthenticateUser authenticates a user
func (m *postgresDBRepo) AuthenticateUser(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var userID int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "SELECT id, password FROM users WHERE email = $1", email)
	err := row.Scan(&userID, &hashedPassword)
	if err != nil {
		return userID, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return userID, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}
	return userID, hashedPassword, nil
}
