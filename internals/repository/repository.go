package repository

import "bookings/internals/models"

type DatabaseRepo interface {
	AllUser() bool
	InsertReservation(models.Reservation) error
}
