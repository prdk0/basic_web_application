package dbrepo

import (
	"bookings/internals/models"
	"context"
	"fmt"
	"time"
)

func (m *postgreDbRepo) AllUser() bool {
	return true
}

func (m *postgreDbRepo) InsertReservation(res models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	var id int
	defer cancel()
	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at) values($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	row := m.DB.QueryRowContext(ctx, stmt, res.FirstName, res.LastName, res.Email, res.Phone, res.StartDate, res.EndDate, res.RoomID, time.Now(), time.Now())
	err := row.Scan(&id)
	if err != nil {
		return err
	}
	fmt.Println(id)
	return nil
}
