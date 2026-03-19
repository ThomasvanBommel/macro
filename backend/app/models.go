package main

type RegisterUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Session struct {
	User_ID   int    `json:"user_id"`
	Username  string `json:"username"`
	ExpiresAt string `json:"expires_at"`
}

type Food struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Brand     string `json:"brand,omitempty"`
	CreatedAt string `json:"created_at"`
	CreatedBy int    `json:"created_by"`
	Calories  int    `json:"calories"`
	Carbs     int    `json:"carbs"`
	Protein   int    `json:"protein"`
	Fat       int    `json:"fat"`
}