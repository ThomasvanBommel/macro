package main

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Message  string `json:"message"`
}
