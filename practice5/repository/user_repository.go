package repository

import (
	"database/sql"
	
	"practice5/model"
)

type Repository struct {
	DB *sql.DB
}

func (r *Repository) GetPaginatedUsers(page int, pageSize int, gender string, orderBy string) (model.PaginatedResponse, error) {

	offset := (page - 1) * pageSize

	query := `
	SELECT id,name,email,gender,birth_date
	FROM users
	WHERE ($1 = '' OR gender = $1)
	ORDER BY ` + orderBy + `
	LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.Query(query, gender, pageSize, offset)

	if err != nil {
		return model.PaginatedResponse{}, err
	}

	defer rows.Close()

	var users []model.User

	for rows.Next() {

		var u model.User

		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.Gender,
			&u.BirthDate,
		)

		if err != nil {
			return model.PaginatedResponse{}, err
		}

		users = append(users, u)
	}

	var total int

	err = r.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&total)

	if err != nil {
		return model.PaginatedResponse{}, err
	}

	return model.PaginatedResponse{
		Data:       users,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}
func (r *Repository) GetCommonFriends(user1 string, user2 string) ([]model.User, error) {

	query := `
	SELECT u.id, u.name, u.email, u.gender, u.birth_date
	FROM users u
	JOIN user_friends f1 ON u.id = f1.friend_id
	JOIN user_friends f2 ON u.id = f2.friend_id
	WHERE f1.user_id = $1 AND f2.user_id = $2
	`

	rows, err := r.DB.Query(query, user1, user2)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []model.User

	for rows.Next() {

		var u model.User

		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.Gender,
			&u.BirthDate,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}