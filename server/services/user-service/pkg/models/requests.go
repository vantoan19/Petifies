package models

type CreateUserReq struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

type LoginReq struct {
	Email    string
	Password string
}
