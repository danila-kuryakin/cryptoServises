package repository

import (
	"controller/pkg/models"
	controllerRepository "controller/pkg/repository"
	"database/sql"
	"errors"
	"fmt"
	customError "telegramBot/internal/error"
)

type UserPostgres struct {
	db *sql.DB
}

func NewUserPostgres(db *sql.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (p *UserPostgres) GetUserById(userId int64) (*models.User, error) {

	query := fmt.Sprintf(`SELECT * FROM %s WHERE user_id = $1`, controllerRepository.UserTable)

	var user models.User
	err := p.db.QueryRow(query, userId).Scan(&user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customError.ErrorUserNotFound
		}
		return nil, errors.New(fmt.Sprint("Error getting user by id ", userId))
	}
	return &user, nil
}

func (p *UserPostgres) CreateUser(userId int64) error {
	query := fmt.Sprintf(`INSERT INTO %s (user_id, proposals_subscribed, spaces_subscribed) VALUES ($1, $2, $3)`, controllerRepository.UserTable)
	_, err := p.db.Exec(query, userId, 0, 0)

	if err != nil {
		return errors.New(fmt.Sprint("Error creating user ", userId))
	}
	return nil
}

func (p *UserPostgres) SetSubscribedSpaces(userId int64, subscribeStatus int) (bool, error) {
	query := fmt.Sprintf(`
        UPDATE %s SET spaces_subscribed = $1
        WHERE user_id = $2`, controllerRepository.UserTable)

	result, err := p.db.Exec(query, subscribeStatus, userId)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if rows == 0 {
		return false, customError.ErrorUserNotFound
	}
	return true, nil
}

func (p *UserPostgres) SetSubscribedProposals(userId int64, subscribeStatus int) (bool, error) {
	query := fmt.Sprintf(`
        UPDATE %s SET proposals_subscribed = $1
        WHERE user_id = $2`, controllerRepository.UserTable)

	result, err := p.db.Exec(query, subscribeStatus, userId)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if rows == 0 {
		return false, customError.ErrorUserNotFound
	}
	return true, nil
}

func (p *UserPostgres) StatusSubscribedSpaces(userId int64) (int, error) {
	query := fmt.Sprintf(`SELECT spaces_subscribed FROM %s WHERE user_id = $1`, controllerRepository.UserTable)

	var spaceSub int
	err := p.db.QueryRow(query, userId).Scan(&spaceSub)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, customError.ErrorUserNotFound
		}
		return 0, err
	}
	return spaceSub, nil
}

func (p *UserPostgres) StatusSubscribedProposals(userId int64) (int, error) {
	query := fmt.Sprintf(`SELECT proposals_subscribed FROM %s WHERE user_id = $1`, controllerRepository.UserTable)

	var proposalSub int
	err := p.db.QueryRow(query, userId).Scan(&proposalSub)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, customError.ErrorUserNotFound
		}
		return 0, err
	}
	return proposalSub, nil
}

func (p *UserPostgres) CreateVotesId(userId int64, votesId string) (bool, error) {
	query := fmt.Sprintf(`
        INSERT INTO %s (user_id, votes_id)
        VALUES ($1, $2)`, controllerRepository.UserVotesTable)

	result, err := p.db.Exec(query, userId, votesId)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if rows == 0 {
		return false, customError.ErrorUserNotFound
	}
	return true, nil
}

func (p *UserPostgres) DropVotesId(userId int64, votesId string) (bool, error) {
	query := fmt.Sprintf(`
        DELETE FROM %s
		WHERE user_id = $1 AND votes_id = $2`, controllerRepository.UserVotesTable)

	result, err := p.db.Exec(query, userId, votesId)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if rows == 0 {
		return false, customError.ErrorUserNotFound
	}
	return true, nil
}

func (p *UserPostgres) GetVotesByUser(userId int64) ([]string, error) {
	query := fmt.Sprintf(`SELECT COALESCE(votes_id, '{}') FROM %s WHERE user_id = $1`, controllerRepository.UserVotesTable)

	var votesId []string
	rows, err := p.db.Query(query, userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, customError.ErrorNotFound
			}
			return nil, err
		}
		votesId = append(votesId, id)
	}
	return votesId, nil
}
