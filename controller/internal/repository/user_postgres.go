package repository

import (
	"database/sql"
	"fmt"
	"log"
)

type UserPostgres struct {
	db *sql.DB
}

func NewUserPostgres(db *sql.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (u *UserPostgres) GetUserSubscriptions() ([]int64, error) {
	query := fmt.Sprintf(`SELECT user_id			       
				FROM %s 
				WHERE dao_subscribed = 1`, userTable)

	rows, err := u.db.Query(query)
	if err != nil {
		log.Println("Query:", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("Close:", err)
		}
	}(rows)

	var users []int64

	for rows.Next() {
		var u int64
		if err := rows.Scan(&u); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
