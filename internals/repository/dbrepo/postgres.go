package dbrepo

import (
	"bookings/internals/models"
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (m *postgreDbRepo) AllUser() bool {
	return true
}

func (m *postgreDbRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var id int
	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at) values($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	row := m.DB.QueryRowContext(ctx, stmt, res.FirstName, res.LastName, res.Email, res.Phone, res.StartDate, res.EndDate, res.RoomID, time.Now(), time.Now())
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *postgreDbRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id, created_at, updated_at, restriction_id) values($1, $2, $3, $4, $5, $6, $7)`
	_, err := m.DB.ExecContext(ctx, stmt, r.StartDate, r.EndDate, r.RoomID, r.ReservationID, time.Now(), time.Now(), r.RestrictionID)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgreDbRepo) SearchAvailabilityByDatesByRoomId(start, end time.Time, roomId int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var numRows int
	query := `select count(id) from room_restrictions where room_id = $1 and $2 < end_date and $3 > start_date;`
	row := m.DB.QueryRowContext(ctx, query, roomId, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}
	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

func (m *postgreDbRepo) SearchAvailabilityForAllrooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var rooms []models.Room
	query := `SELECT r.id, r.room_name
			  FROM rooms r
			  WHERE r.id NOT IN (SELECT rr.room_id FROM room_restrictions rr WHERE $1 < end_date and $2 > start_date);`
	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}
	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}
	if err := rows.Err(); err != nil {
		return rooms, nil
	}
	return rooms, nil
}

func (m *postgreDbRepo) GetRoomById(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var room models.Room
	query := `select id, room_name from rooms where id = $1`
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&room.ID, &room.RoomName)
	if err != nil {
		return room, err
	}
	return room, nil
}

func (m *postgreDbRepo) GetRestrictions() ([]models.Restriction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var restrictions []models.Restriction
	query := `select id, restriction_name from restrictions;`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return restrictions, err
	}
	for rows.Next() {
		var restriction models.Restriction
		err := rows.Scan(&restriction.ID, &restriction.RestrictionName)
		if err != nil {
			return restrictions, err
		}
		restrictions = append(restrictions, restriction)
	}

	if err := rows.Err(); err != nil {
		return restrictions, err
	}

	return restrictions, nil
}

func (m *postgreDbRepo) GetUserById(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `select id, first_name, last_name, email, password, access_level, created_at, updated_at
			from users where id=$1`
	row := m.DB.QueryRowContext(ctx, query, id)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.AccessLevel,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (m *postgreDbRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `update users set first_name=$1, last_name=$2, email=$3, access_level=$4, updated_at = $5`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *postgreDbRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select id, password from users where email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

func (m *postgreDbRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var reservations []models.Reservation
	query := `select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at, r.processed, rm.id, rm.room_name from reservations r left join rooms rm on (r.room_id = rm.id) order by r.start_date asc`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()
	for rows.Next() {
		var reservation models.Reservation
		err := rows.Scan(
			&reservation.ID,
			&reservation.FirstName,
			&reservation.LastName,
			&reservation.Email,
			&reservation.Phone,
			&reservation.StartDate,
			&reservation.EndDate,
			&reservation.RoomID,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
			&reservation.Processed,
			&reservation.Room.ID,
			&reservation.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, reservation)
	}

	if err := rows.Err(); err != nil {
		return reservations, err
	}
	return reservations, nil
}

func (m *postgreDbRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var reservations []models.Reservation
	query := `select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at, rm.id, rm.room_name from reservations r left join rooms rm on (r.room_id = rm.id) 
	where processed = 0
	order by r.start_date asc`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()
	for rows.Next() {
		var reservation models.Reservation
		err := rows.Scan(
			&reservation.ID,
			&reservation.FirstName,
			&reservation.LastName,
			&reservation.Email,
			&reservation.Phone,
			&reservation.StartDate,
			&reservation.EndDate,
			&reservation.RoomID,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
			&reservation.Room.ID,
			&reservation.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, reservation)
	}

	if err := rows.Err(); err != nil {
		return reservations, err
	}
	return reservations, nil
}

func (m *postgreDbRepo) GetReservationById(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var res models.Reservation

	query := `select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at, r.processed, rm.id, rm.room_name 
	from reservations r
	 left join rooms rm on (r.room_id = rm.id)
	 where r.id = $1
	`
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Phone,
		&res.StartDate,
		&res.EndDate,
		&res.RoomID,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.Processed,
		&res.Room.ID,
		&res.Room.RoomName,
	)

	if err != nil {
		return res, err
	}

	return res, nil
}
