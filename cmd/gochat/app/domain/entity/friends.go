package entity

type Friend struct {
	ID       int64 `db:"id"        json:"id"`
	UserID   int64 `db:"user_id"   json:"user_id"`
	FriendID int64 `db:"friend_id" json:"friend_id"`
}
