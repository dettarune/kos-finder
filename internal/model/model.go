package model


type User struct {
	Username string
}

type UserLogin struct {
	Username             string `json:"username"`
	Password             string `json:"password"`
}