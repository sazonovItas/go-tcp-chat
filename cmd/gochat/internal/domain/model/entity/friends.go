package entity

type Friend struct {
	ID       int64 `db:"id"`
	UserID   int64 `db:"user_id"`
	FriendID int64 `db:"friend_id"`
}
