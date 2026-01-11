package models

type User struct {
	Id            int64 `db:"id"`
	UserId        int64 `json:"user_id"`
	DaoSubscribed int   `json:"dao_subscribed"`
}
