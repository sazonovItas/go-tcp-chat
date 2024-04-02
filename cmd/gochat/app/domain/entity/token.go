package entity

type TokenID string

func (tid TokenID) String() string {
	return string(tid)
}

type Token struct {
	ID           TokenID `json:"id"`
	UserId       int64   `json:"user_id"`
	Login        string  `json:"login"`
	PasswordHash string  `json:"password_hash"`
}
