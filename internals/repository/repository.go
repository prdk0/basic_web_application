package repository

import (
	"bookings/internals/models"
	"time"
)

type DatabaseRepo interface {
	AllUser() bool
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomId(start, end time.Time, roomId int) (bool, error)
	SearchAvailabilityForAllrooms(start, end time.Time) ([]models.Room, error)
}
