package repository

import (
	"time"

	"github.com/p3rfect05/go_proj/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityByDateByRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	SearchRoomByID(roomID int) (models.Room, error)
	GetUserByID(ID int) (models.User, error)
	AuthenticateUser(email, testPassword string) (int, string, error)
	UpdateUser(u models.User) error
}
