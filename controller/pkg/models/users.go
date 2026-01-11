package models

type User struct {
	Id            int64 `json:"id" db:"id"`
	UserId        int64 `json:"user_id" db:"user_id"`
	DaoSubscribed int   `json:"dao_subscribed" db:"dao_subscribed"`
}
