package models

type User struct {
	ID       int    `db:"id"`
	Login    string `db:"login"`
	Password []byte `db:"pass_hash"`
}
