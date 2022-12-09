package postgres

import (
	"blockpost/genprotos/authorization"
	"errors"
	"time"
)

// CreateUser ...
func (p Postgres) CreateUser(id string, req *authorization.CreateUserRequest) error {

	_, err := p.DB.Exec(`Insert into "user"("id", "fname", "lname", "username", "password",
	"user_type", "address", "phone", "created_at") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, now())
	`, id, req.Fname, req.Lname, req.Username, req.Password, req.UserType, req.Address, req.Phone)
	if err != nil {
		return err
	}
	return nil
}

// GetUserByID ...
func (p Postgres) GetUserByID(id string) (*authorization.User, error) {
	res := &authorization.User{}
	var deletedAt *time.Time
	var updatedAt *string
	err := p.DB.QueryRow(`SELECT 
							"id", "fname", "lname",	"username",	"password",	"user_type", 
							"address", "phone", "created_at", "updated_at", "deleted_at"
    	FROM "user" WHERE id = $1`, id).Scan(
		&res.Id, &res.Fname, &res.Lname, &res.Username, &res.Password, &res.UserType,
		&res.Address, &res.Phone, &res.CreatedAt, &updatedAt, &deletedAt)
	if err != nil {
		return nil, err
	}

	if updatedAt != nil {
		res.UpdatedAt = *updatedAt
	}

	if deletedAt != nil {
		return res, errors.New("user not found")
	}

	return res, err
}

// GetUserList ...
func (p Postgres) GetUserList(offset, limit int, search string) (*authorization.GetUserListResponse, error) {
	resp := &authorization.GetUserListResponse{
		Users: make([]*authorization.User, 0),
	}
	rows, err := p.DB.Queryx(`SELECT
	"id", "fname", "lname",	"username",	"password",	"user_type", 
	"address", "phone", "created_at", "updated_at"
	FROM "user" WHERE deleted_at IS NULL AND ("username" || "fname" || "lname" ILIKE '%' || $1 || '%')
	LIMIT $2
	OFFSET $3
	`, search, limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		a := &authorization.User{}

		var updatedAt *string

		err := rows.Scan(
			&a.Id, &a.Fname, &a.Lname, &a.Username, &a.Password, &a.UserType,
			&a.Address, &a.Phone, &a.CreatedAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		if updatedAt != nil {
			a.UpdatedAt = *updatedAt
		}

		resp.Users = append(resp.Users, a)
	}

	return resp, err
}

// UpdateUser ...
func (p Postgres) UpdateUser(entity *authorization.UpdateUserRequest) error {

	res, err := p.DB.NamedExec(`UPDATE "user" 
		SET "fname"=:f, "lname"=:l, "username"=:u, "password"=:p, "user_type"=:ut, "address"=:a, "phone"=:ph, "updated_at"=now() 
		WHERE deleted_at IS NULL AND id=:id`, map[string]interface{}{
		"id": entity.Id,
		"f":  entity.Fname, "l": entity.Lname, "u": entity.Username, "p": entity.Password,
		"ut": entity.UserType, "a": entity.Address, "ph": entity.Phone,
	})
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n > 0 {
		return nil
	}

	return errors.New("user not found")
}

// DeleteUser ...
func (p Postgres) DeleteUser(id string) error {
	res, err := p.DB.Exec(`UPDATE "user" SET deleted_at=now() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n > 0 {
		return nil
	}

	return errors.New("user not found")
}

// GetUserByUsername ...
func (p Postgres) GetUserByUsername(username string) (*authorization.User, error) {
	res := &authorization.User{}
	var deletedAt *time.Time
	var updatedAt *string
	err := p.DB.QueryRow(`SELECT 
							"id", "fname", "lname",	"username",	"password",	"user_type", 
							"address", "phone", "created_at", "updated_at", "deleted_at"
    	FROM "user" WHERE "username" = $1`, username).Scan(
		&res.Id, &res.Fname, &res.Lname, &res.Username, &res.Password, &res.UserType,
		&res.Address, &res.Phone, &res.CreatedAt, &updatedAt, &deletedAt)
	if err != nil {
		return nil, err
	}

	if updatedAt != nil {
		res.UpdatedAt = *updatedAt
	}

	if deletedAt != nil {
		return res, errors.New("user not found")
	}

	return res, err
}
