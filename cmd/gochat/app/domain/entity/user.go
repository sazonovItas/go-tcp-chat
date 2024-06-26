package entity

type User struct {
	ID           int64  `db:"id"            json:"id"`
	Login        string `db:"login"         json:"login"`
	Name         string `db:"name"          json:"name"`
	Color        string `db:"color"         json:"color"`
	PasswordHash string `db:"password_hash" json:"password_hash"`
}

type PublicUser struct {
	ID    int64  `db:"id"    json:"id"`
	Login string `db:"login" json:"login"`
	Name  string `db:"name"  json:"name"`
	Color string `db:"color" json:"color"`
}

type AuthUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
