package main

type RegisterUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterUserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Message  string `json:"message"`
}
