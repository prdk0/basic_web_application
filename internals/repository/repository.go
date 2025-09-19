package repository

import "bookings/internals/models"

type DatabaseRepo interface {
	AllUser() bool
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
}
